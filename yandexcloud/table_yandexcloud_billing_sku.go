package yandexcloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableYandexBillingSku(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_billing_sku",
		Description: "Yandex Cloud Billing SKUs (service catalog).",
		List: &plugin.ListConfig{
			Hydrate: listYandexBillingSkus,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getYandexBillingSku,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "SKU ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "SKU name."},
			{Name: "service_id", Type: proto.ColumnType_STRING, Description: "Service ID."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "SKU description."},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "SKU currency."},
		},
	}
}

func listYandexBillingSkus(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return nil, nil
}

func getYandexBillingSku(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return nil, nil
}
