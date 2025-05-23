package yandexcloud

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/context_key"
)

var currentLogLevel LogLevel = LogLevelError

// LogLevel defines the logging level
// LogLevel is a type for logging levels
// (add this definition if it is not defined elsewhere)
//
//go:generate stringer -type=LogLevel

func parseLogLevel(val interface{}) LogLevel {
	switch v := val.(type) {
	case LogLevel:
		return v
	case *LogLevel:
		if v != nil {
			return *v
		}
	case string:
		switch v {
		case "debug", "DEBUG":
			return LogLevelDebug
		case "info", "INFO":
			return LogLevelInfo
		case "error", "ERROR":
			return LogLevelError
		}
	case *string:
		if v != nil {
			switch *v {
			case "debug", "DEBUG":
				return LogLevelDebug
			case "info", "INFO":
				return LogLevelInfo
			case "error", "ERROR":
				return LogLevelError
			}
		}
	}
	return LogLevelError // default
}

func SetLogLevelFromConfig(cfg *Config) {
	if cfg == nil || cfg.LogLevel == nil {
		LogError(context.Background(), "SetLogLevelFromConfig: log level is nil")
		return
	}
	level := parseLogLevel(cfg.LogLevel)
	LogInfo(context.Background(), "SetLogLevelFromConfig: setting log level to %v", level)
	currentLogLevel = level
}

func ShouldLog(level LogLevel) bool {
	switch currentLogLevel {
	case LogLevelDebug:
		return true
	case LogLevelInfo:
		return level != LogLevelDebug
	case LogLevelError:
		return level == LogLevelError
	default:
		return false
	}
}

func init() {
	log.SetOutput(os.Stdout)
}

// Strict error type for authentication
type AuthError string

func (e AuthError) Error() string { return string(e) }

// getConfig safely returns *Config from connection config handling both pointer and value cases.
func getConfig(d *plugin.QueryData) *Config {
	LogInfo(context.Background(), "getConfig called")
	if d == nil || d.Connection == nil || d.Connection.Config == nil {
		LogDebug(context.Background(), "getConfig: d, d.Connection or d.Connection.Config == nil")
		return &Config{}
	}
	var c *Config
	if v, ok := d.Connection.Config.(*Config); ok {
		c = v
	} else if v, ok := d.Connection.Config.(Config); ok {
		c = &v
	} else {
		LogDebug(context.Background(), "getConfig: d.Connection.Config is not *Config or Config")
		return &Config{}
	}
	SetLogLevelFromConfig(c)
	// Convert string values to strict types if needed
	if c.Token != nil {
		t := Token(*((*string)(c.Token)))
		c.Token = &t
	}
	if c.CloudID != nil {
		cid := CloudID(*((*string)(c.CloudID)))
		c.CloudID = &cid
	}
	if c.FolderID != nil {
		fid := FolderID(*((*string)(c.FolderID)))
		c.FolderID = &fid
	}
	if c.UserAgent != nil {
		u := UserAgent(*((*string)(c.UserAgent)))
		c.UserAgent = &u
	}
	if c.EndpointOverride != nil {
		e := EndpointOverride(*((*string)(c.EndpointOverride)))
		c.EndpointOverride = &e
	}
	LogDebug(context.Background(), "getConfig: config = Token=%v, ServiceAccountKeyFile=%v, CloudID=%v, FolderID=%v", c.Token, derefString(c.ServiceAccountKeyFile), derefString(c.CloudID), derefString(c.FolderID))
	if err := ValidateConfig(c); err != nil {
		LogError(context.Background(), "Config validation failed: %v", err)
		return &Config{}
	}
	return c
}

// getQualString extracts string qual or returns default.
func getQualString(d *plugin.QueryData, column string, defaultValue *string) string {
	if q, ok := d.KeyColumnQuals[column]; ok {
		if s := q.GetStringValue(); s != "" {
			return s
		}
	}
	if defaultValue != nil {
		return *defaultValue
	}
	return ""
}

// httpClientCache caches http.Client by timeout.
var httpClientCache sync.Map // map[int64]*http.Client

// GetHTTPClient returns http.Client with the specified timeout (seconds).
func GetHTTPClient(timeoutSec int64) *http.Client {
	if timeoutSec <= 0 {
		timeoutSec = 30 // default
	}
	if v, ok := httpClientCache.Load(timeoutSec); ok {
		return v.(*http.Client)
	}
	c := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	httpClientCache.Store(timeoutSec, c)
	return c
}

// getLoggerFromContext safely retrieves the logger from context or returns nil if not present.
func getLoggerFromContext(ctx context.Context) hclog.Logger {
	logger, ok := ctx.Value(context_key.Logger).(hclog.Logger)
	if !ok {
		return nil
	}
	return logger
}

// LogDebug logs debug messages, falling back to log.Printf if logger is missing.
func LogDebug(ctx context.Context, format string, args ...interface{}) {
	if logger := getLoggerFromContext(ctx); logger != nil {
		logger.Debug(fmt.Sprintf(format, args...))
	} else {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// LogInfo logs info messages, falling back to log.Printf if logger is missing.
func LogInfo(ctx context.Context, format string, args ...interface{}) {
	if logger := getLoggerFromContext(ctx); logger != nil {
		logger.Info(fmt.Sprintf(format, args...))
	} else {
		log.Printf("[INFO] "+format, args...)
	}
}

// LogError logs error messages, falling back to log.Printf if logger is missing.
func LogError(ctx context.Context, format string, args ...interface{}) {
	if logger := getLoggerFromContext(ctx); logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	} else {
		log.Printf("[ERROR] "+format, args...)
	}
}

// HandleHTTPError centrally handles HTTP response errors.
// ignoreCodes - list of codes that are considered not an error (e.g., 404, 403).
func HandleHTTPError(resp *http.Response, ignoreCodes ...int) error {
	var urlStr string
	if resp.Request != nil && resp.Request.URL != nil {
		urlStr = resp.Request.URL.String()
	}
	for _, code := range ignoreCodes {
		if resp.StatusCode == code {
			log.Printf("HTTP %d ignored for %s", resp.StatusCode, urlStr)
			return nil
		}
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	var body []byte
	if resp.Body != nil {
		body, _ = io.ReadAll(resp.Body)
	}
	log.Printf("API error %d for %s: %s", resp.StatusCode, urlStr, string(body))
	return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
}

// DoWithRetry performs HTTP request with retry support for temporary errors (5xx, 429).
// ctx - context for logging and cancellation, retryCount - number of attempts (>=1), timeoutSec - timeout for each attempt.
// reqFactory - function returning new *http.Request for each attempt.
func DoWithRetry(ctx context.Context, client *http.Client, reqFactory func() *http.Request, retryCount int, timeoutSec int64) (*http.Response, error) {
	if retryCount < 1 {
		retryCount = 1
	}
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	var lastErr error
	for attempt := 1; attempt <= retryCount; attempt++ {
		req := reqFactory()
		cctx, cancel := context.WithTimeout(req.Context(), time.Duration(timeoutSec)*time.Second)
		req = req.Clone(cctx)
		resp, err := client.Do(req)
		if err != nil {
			LogError(ctx, "HTTP request failed (attempt %d/%d) for %s: %v", attempt, retryCount, req.URL, err)
			lastErr = err
			cancel()
			continue
		}
		if resp.StatusCode < 500 && resp.StatusCode != 429 {
			cancel()
			return resp, nil
		}
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		resp.Body.Close()
		LogError(ctx, "HTTP %d (attempt %d/%d) for %s: %s", resp.StatusCode, attempt, retryCount, req.URL, string(body))
		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		cancel()
		time.Sleep(time.Duration(attempt*200) * time.Millisecond)
	}
	return nil, lastErr
}

// ApplyRequestOptions applies user-agent and endpoint override to http.Request.
func ApplyRequestOptions(req *http.Request, userAgent, endpointOverride *string) {
	if userAgent != nil && *userAgent != "" {
		req.Header.Set("User-Agent", *userAgent)
	}
	if endpointOverride != nil && *endpointOverride != "" {
		// Only change scheme+host, keep path/query
		u := *req.URL
		override, err := url.Parse(*endpointOverride)
		if err == nil {
			u.Scheme = override.Scheme
			u.Host = override.Host
			req.URL = &u
		}
	}
}

func derefString(ptr interface{}) string {
	switch v := ptr.(type) {
	case *string:
		if v != nil {
			return *v
		}
	case *CloudID:
		if v != nil {
			return string(*v)
		}
	case *FolderID:
		if v != nil {
			return string(*v)
		}
	}
	return "<nil>"
}
