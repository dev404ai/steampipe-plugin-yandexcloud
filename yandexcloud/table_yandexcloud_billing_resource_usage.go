package yandexcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// Strict types for identifiers and parameters
// type BillingAccountID string
// type BillingAccountName string

// BillingAccount carries minimal YC billing account info.
// type BillingAccount struct {
// 	Id   BillingAccountID   `json:"id"`
// 	Name BillingAccountName `json:"name"`
// }

type BillingAccountsResponse struct {
	BillingAccounts []BillingAccount `json:"billingAccounts"`
}

func tableYandexBillingResourceUsage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_billing_resource_usage",
		Description: "Billing account usage & cost records.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"billing_account_id", "cloud_id", "folder_id"}),
			Hydrate:    listYandexBillingUsage,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("billing_account_id"),
			Hydrate:    getYandexBillingAccountResourceUsage,
		},
		Columns: []*plugin.Column{
			{Name: "billing_account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Billing account ID."},
			{Name: "billing_account_name", Type: proto.ColumnType_STRING, Transform: transform.From(billingAccountNameUpperTransform), Description: "Billing account name (UPPERCASE)."},
			{Name: "cloud_id", Type: proto.ColumnType_STRING, Description: "Cloud ID that owns the billing account."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Description: "Folder ID associated with the record."},
		},
	}
}

func listYandexBillingUsage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	_ = ctx // suppress unused parameter warning
	cfg := getConfig(d)
	token, errTok := getAuthToken(ctx, cfg)
	if errTok != nil {
		LogError(ctx, "Billing: failed to get auth token: %v", errTok)
		return nil, errTok
	}

	endpoint := "https://billing.api.cloud.yandex.net/billing/v1/billingAccounts"
	LogInfo(ctx, "Billing: GET %s", endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		LogError(ctx, "Billing: failed to create request: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
	timeoutSec := TimeoutSec(30)
	if cfg.Timeout != nil && *cfg.Timeout > 0 {
		timeoutSec = TimeoutSec(*cfg.Timeout)
	}
	retryCount := RetryCount(3)
	if cfg.Retry != nil && *cfg.Retry > 0 {
		retryCount = RetryCount(*cfg.Retry)
	}
	client := GetHTTPClient(int64(timeoutSec))
	resp, err := DoWithRetry(ctx, client, func() *http.Request { return req }, int(retryCount), int64(timeoutSec))
	if err != nil {
		LogError(ctx, "Billing: request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if err := HandleHTTPError(resp); err != nil {
		LogError(ctx, "Billing: HTTP error: %v", err)
		return nil, err
	}
	var result BillingAccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		LogError(ctx, "Billing: failed to decode response: %v", err)
		return nil, err
	}
	LogInfo(ctx, "Billing: success, got %d accounts", len(result.BillingAccounts))
	for _, acc := range result.BillingAccounts {
		yield := map[string]interface{}{
			"billing_account_id":   acc.Id,
			"billing_account_name": acc.Name,
			"_ctx":                 ctx,
			"_token":               token,
		}
		return yield, nil // For example, return the first account
	}
	return nil, nil
}

// Hydrate function to get a single account by billing_account_id
func getYandexBillingAccountResourceUsage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var accountID string
	if h != nil && h.Item != nil {
		if acc, ok := h.Item.(*BillingAccount); ok {
			accountID = string(acc.Id)
		}
	}
	if accountID == "" {
		if v, ok := d.KeyColumnQuals["billing_account_id"]; ok {
			accountID = v.GetStringValue()
		}
	}
	if accountID == "" {
		return nil, fmt.Errorf("billing_account_id must be provided")
	}
	cfg := getConfig(d)
	token, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	endpoint := "https://billing.api.cloud.yandex.net/billing/v1/billingAccounts/" + accountID
	LogInfo(ctx, "Billing: GET %s", endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
	timeoutSec := TimeoutSec(30)
	if cfg.Timeout != nil && *cfg.Timeout > 0 {
		timeoutSec = TimeoutSec(*cfg.Timeout)
	}
	retryCount := RetryCount(3)
	if cfg.Retry != nil && *cfg.Retry > 0 {
		retryCount = RetryCount(*cfg.Retry)
	}
	client := GetHTTPClient(int64(timeoutSec))
	resp, err := DoWithRetry(ctx, client, func() *http.Request { return req }, int(retryCount), int64(timeoutSec))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := HandleHTTPError(resp); err != nil {
		return nil, err
	}
	var acc BillingAccount
	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

// Transform function: account name in uppercase
func billingAccountNameUpperTransform(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	acc, ok := d.HydrateItem.(*BillingAccount)
	if !ok || acc.Name == "" {
		return nil, nil
	}
	return strings.ToUpper(string(acc.Name)), nil
}
