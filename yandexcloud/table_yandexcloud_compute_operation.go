package yandexcloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeOperation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_operation",
		Description: "Yandex Cloud Compute operation (get by ID only).",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("operation_id"),
			Hydrate:    getYandexComputeOperation,
		},
		Columns: []*plugin.Column{
			{Name: "operation_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Operation ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Operation status."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Operation description."},
			{Name: "done", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Done"), Description: "Operation done flag."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtOperationDateTransform), Description: "Operation creation date (YYYY-MM-DD)."},
			{Name: "error", Type: proto.ColumnType_JSON, Transform: transform.FromField("Error"), Description: "Operation error details."},
			{Name: "response", Type: proto.ColumnType_JSON, Transform: transform.FromField("Response"), Description: "Operation response details."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Transform: transform.FromField("Metadata"), Description: "Operation metadata."},
		},
	}
}

func getYandexComputeOperation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var opID string
	if h != nil && h.Item != nil {
		if op, ok := h.Item.(*Operation); ok {
			opID = op.Id
		}
	}
	if opID == "" {
		if v, ok := d.KeyColumnQuals["operation_id"]; ok {
			opID = v.GetStringValue()
		}
	}
	if opID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	op, err := client.GetOperation(ctx, OperationID(opID), 30, 3)
	if err != nil {
		return nil, err
	}
	return op, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtOperationDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	op, ok := d.HydrateItem.(*Operation)
	if !ok || op.CreatedAt == "" {
		return nil, nil
	}
	if len(op.CreatedAt) < 10 {
		return op.CreatedAt, nil
	}
	return op.CreatedAt[:10], nil
}
