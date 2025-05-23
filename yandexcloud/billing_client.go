package yandexcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// There is no BillingAccount struct in this file; it is defined in table_yandexcloud_billing_account.go

type ListBillingAccountsResponse struct {
	BillingAccounts []*BillingAccount `json:"billingAccounts"`
	NextPageToken   string            `json:"nextPageToken"`
}

type GetBillingAccountResponse struct {
	BillingAccount *BillingAccount `json:"billingAccount"`
}

type BillingClient interface {
	ListBillingAccounts(ctx context.Context, token, pageToken string, pageSize int64, timeoutSec int64) ([]*BillingAccount, string, error)
	GetBillingAccount(ctx context.Context, token, accountID string, timeoutSec int64) (*BillingAccount, error)
}

type yandexBillingClient struct{}

func NewBillingClient() BillingClient {
	return &yandexBillingClient{}
}

func (c *yandexBillingClient) ListBillingAccounts(ctx context.Context, token, pageToken string, pageSize int64, timeoutSec int64) ([]*BillingAccount, string, error) {
	const endpoint = "https://billing.api.cloud.yandex.net/billing/v1/billingAccounts"
	params := url.Values{}
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if pageSize > 0 {
		params.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	urlStr := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	LogInfo(ctx, "BillingClient: GET %s", urlStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		LogError(ctx, "BillingClient: failed to create request: %v", err)
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		LogError(ctx, "BillingClient: request failed: %v", err)
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		LogError(ctx, "BillingClient: API error: %s", resp.Status)
		return nil, "", fmt.Errorf("billing API error: %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	LogInfo(ctx, "BillingClient: raw response: %s", string(body))
	var respBody ListBillingAccountsResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		LogError(ctx, "BillingClient: failed to decode response: %v", err)
		return nil, "", err
	}
	LogInfo(ctx, "BillingClient: got %d accounts", len(respBody.BillingAccounts))
	return respBody.BillingAccounts, respBody.NextPageToken, nil
}

func (c *yandexBillingClient) GetBillingAccount(ctx context.Context, token, accountID string, timeoutSec int64) (*BillingAccount, error) {
	urlStr := fmt.Sprintf("https://billing.api.cloud.yandex.net/billing/v1/billingAccounts/%s", accountID)
	LogInfo(ctx, "BillingClient: GET %s", urlStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		LogError(ctx, "BillingClient: failed to create request: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		LogError(ctx, "BillingClient: request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		LogError(ctx, "BillingClient: API error: %s", resp.Status)
		return nil, fmt.Errorf("billing API error: %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	var respBody GetBillingAccountResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		LogError(ctx, "BillingClient: failed to decode response: %v", err)
		return nil, err
	}
	LogInfo(ctx, "BillingClient: got account %s", accountID)
	return respBody.BillingAccount, nil
}
