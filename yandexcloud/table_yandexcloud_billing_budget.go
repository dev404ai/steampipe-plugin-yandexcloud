package yandexcloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableYandexBillingBudget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_billing_budget",
		Description: "Yandex Cloud Billing Budgets.",
		List: &plugin.ListConfig{
			Hydrate: listYandexBillingBudgets,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getYandexBillingBudget,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Budget ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Budget name."},
			{Name: "billing_account_id", Type: proto.ColumnType_STRING, Description: "Billing account ID."},
			{Name: "amount", Type: proto.ColumnType_STRING, Description: "Budget amount."},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "Budget currency."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Budget status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Description: "Budget creation date (RFC3339)."},
		},
	}
}

func listYandexBillingBudgets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return nil, nil
}

func getYandexBillingBudget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return nil, nil
}
