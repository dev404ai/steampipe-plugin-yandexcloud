package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableYandexBillingAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_billing_account",
		Description: "Yandex Cloud Billing Accounts.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"id", "name", "active"}),
			Hydrate:    listYandexBillingAccounts,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getYandexBillingAccount,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Billing account ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Account name."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Description: "Creation date (RFC3339)."},
			{Name: "country_code", Type: proto.ColumnType_STRING, Description: "Country code."},
			{Name: "balance", Type: proto.ColumnType_STRING, Description: "Current balance."},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "Currency code."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Is account active?"},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexBillingAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := GetConfig(d.Connection)
	token, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	var filters []string
	if id := getQualString(d, "id", nil); id != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", id))
	}
	if n := getQualString(d, "name", nil); n != "" {
		filters = append(filters, fmt.Sprintf("(name = \"%s\")", n))
	}
	if a, ok := d.KeyColumnQuals["active"]; ok {
		if a.GetBoolValue() {
			filters = append(filters, "(active = true)")
		} else {
			filters = append(filters, "(active = false)")
		}
	}
	pageToken := ""
	pageSize := int64(1000)
	timeoutSec := int64(30)
	if cfg.Timeout != nil && *cfg.Timeout > 0 {
		timeoutSec = int64(*cfg.Timeout)
	}
	client := NewBillingClient()
	for {
		accounts, nextPageToken, err := client.ListBillingAccounts(ctx, token, pageToken, pageSize, timeoutSec)
		if err != nil {
			return nil, err
		}
		for _, acc := range accounts {
			if len(filters) > 0 {
				if !billingAccountMatchesFilters(acc, filters) {
					continue
				}
			}
			if acc != nil {
				d.StreamListItem(ctx, map[string]interface{}{
					"id":           acc.Id,
					"name":         acc.Name,
					"created_at":   acc.CreatedAt,
					"country_code": acc.CountryCode,
					"balance":      acc.Balance,
					"currency":     acc.Currency,
					"active":       acc.Active,
					"labels":       acc.Labels,
				})
			}
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexBillingAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cfg := GetConfig(d.Connection)
	token, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	var id string
	if h != nil && h.Item != nil {
		if acc, ok := h.Item.(*BillingAccount); ok {
			id = acc.Id
		}
	}
	if id == "" {
		if v, ok := d.KeyColumnQuals["id"]; ok {
			id = v.GetStringValue()
		}
	}
	if id == "" {
		return nil, nil // id is missing, not a token error
	}
	timeoutSec := int64(30)
	if cfg.Timeout != nil && *cfg.Timeout > 0 {
		timeoutSec = int64(*cfg.Timeout)
	}
	client := NewBillingClient()
	acc, err := client.GetBillingAccount(ctx, token, id, timeoutSec)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// Manual filter for Billing Account, since API does not support all filters
func billingAccountMatchesFilters(acc *BillingAccount, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, acc.Id) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, acc.Name) {
			return false
		}
		if f == "(active = true)" && !acc.Active {
			return false
		}
		if f == "(active = false)" && acc.Active {
			return false
		}
	}
	return true
}
