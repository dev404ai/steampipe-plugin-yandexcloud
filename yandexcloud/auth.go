package yandexcloud

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type cachedToken struct {
	Token     string
	ExpiresAt time.Time
}

var tokenCache sync.Map

// getAuthToken returns OAuth token or exchanges SA key JSON for IAM token.
func getAuthToken(ctx context.Context, cfg *Config) (string, error) {
	LogInfo(ctx, "getAuthToken called")
	LogDebug(ctx, "getAuthToken: config = %+v", cfg)
	if err := ValidateConfig(cfg); err != nil {
		LogError(ctx, "Config validation failed: %v", err)
		return "", AuthError(err.Error())
	}
	if cfg == nil {
		LogError(ctx, "Config is nil in getAuthToken")
		return "", AuthError("config is nil")
	}
	if cfg.Token != nil && *cfg.Token != "" {
		LogInfo(ctx, "Using static token for authentication")
		return string(*cfg.Token), nil
	}
	if cfg.ServiceAccountKeyFile == nil || *cfg.ServiceAccountKeyFile == "" {
		LogError(ctx, "Neither token nor service_account_key_file set in config")
		return "", AuthError("token or service_account_key_file must be set in connection config")
	}
	path := *cfg.ServiceAccountKeyFile

	// cached?
	if v, ok := tokenCache.Load(path); ok {
		ct := v.(cachedToken)
		if time.Until(ct.ExpiresAt) > 5*time.Minute {
			LogInfo(ctx, "Using cached IAM token for %s", path)
			return ct.Token, nil
		}
	}

	// read key file
	data, err := os.ReadFile(path)
	if err != nil {
		LogError(ctx, "Failed to read key file %s: %v", path, err)
		return "", AuthError("read key file: " + err.Error())
	}
	var keyFile struct {
		ID               string `json:"id"`
		ServiceAccountID string `json:"service_account_id"`
		PrivateKey       string `json:"private_key"`
	}
	if err := json.Unmarshal(data, &keyFile); err != nil {
		LogError(ctx, "Failed to parse key JSON: %v", err)
		return "", AuthError("parse key json: " + err.Error())
	}

	block, _ := pem.Decode([]byte(keyFile.PrivateKey))
	if block == nil {
		LogError(ctx, "Invalid PEM in private_key")
		return "", AuthError("invalid PEM in private_key")
	}
	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		LogError(ctx, "Failed to parse pkcs8 key: %v", err)
		return "", AuthError("parse pkcs8 key: " + err.Error())
	}
	rsaKey, ok := pk.(*rsa.PrivateKey)
	if !ok {
		LogError(ctx, "Key is not RSA")
		return "", AuthError("key is not RSA")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"aud": "https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"iss": keyFile.ServiceAccountID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = keyFile.ID

	jwtStr, err := token.SignedString(rsaKey)
	if err != nil {
		LogError(ctx, "Failed to sign JWT: %v", err)
		return "", AuthError("sign jwt: " + err.Error())
	}

	payload := fmt.Sprintf(`{"jwt":"%s"}`, jwtStr)
	client := GetHTTPClient(30) // default 30s timeout
	req, err := http.NewRequest("POST", "https://iam.api.cloud.yandex.net/iam/v1/tokens", strings.NewReader(payload))
	if err != nil {
		LogError(ctx, "Failed to create IAM request: %v", err)
		return "", AuthError("iam request: " + err.Error())
	}
	var ua *string
	if cfg.UserAgent != nil {
		s := string(*cfg.UserAgent)
		ua = &s
	}
	var eo *string
	if cfg.EndpointOverride != nil {
		s := string(*cfg.EndpointOverride)
		eo = &s
	}
	ApplyRequestOptions(req, ua, eo)
	resp, err := client.Do(req)
	if err != nil {
		LogError(ctx, "IAM request failed: %v", err)
		return "", AuthError("iam request: " + err.Error())
	}
	defer resp.Body.Close()

	if err := HandleHTTPError(resp); err != nil {
		LogError(ctx, "IAM HTTP error: %v", err)
		return "", AuthError("iam http error: " + err.Error())
	}
	var res struct {
		IAMToken  string `json:"iamToken"`
		ExpiresAt string `json:"expiresAt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		LogError(ctx, "Failed to decode IAM response: %v", err)
		return "", AuthError("decode iam response: " + err.Error())
	}
	exp, _ := time.Parse(time.RFC3339, res.ExpiresAt)
	tokenCache.Store(path, cachedToken{Token: res.IAMToken, ExpiresAt: exp})

	LogInfo(ctx, "Successfully obtained new IAM token for %s, expires at %s", path, res.ExpiresAt)
	return res.IAMToken, nil
}
