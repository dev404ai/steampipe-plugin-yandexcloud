package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCOperation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_operation",
		Description: "Yandex Cloud VPC operations.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"operation_id", "description", "created_by", "done"}),
			Hydrate:    listYandexVPCOperations,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("operation_id"),
			Hydrate:    getYandexVPCOperation,
		},
		Columns: []*plugin.Column{
			{Name: "operation_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Operation ID."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Operation description."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCOperationDateTransform), Description: "Operation creation date (YYYY-MM-DD)."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Transform: transform.FromField("CreatedBy"), Description: "ID of the user or service account who initiated the operation."},
			{Name: "modified_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("ModifiedAt"), Description: "The time when the operation was last modified."},
			{Name: "done", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Done"), Description: "If true, the operation is completed."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Transform: transform.FromField("Metadata"), Description: "Service-specific metadata associated with the operation."},
			{Name: "error", Type: proto.ColumnType_JSON, Transform: transform.FromField("Error"), Description: "The error result of the operation in case of failure or cancellation."},
			{Name: "response", Type: proto.ColumnType_JSON, Transform: transform.FromField("Response"), Description: "The normal response of the operation in case of success."},
		},
	}
}

func listYandexVPCOperations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)

	var filters []string
	if id := getQualString(d, "operation_id", nil); id != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", id))
	}
	if dsc := getQualString(d, "description", nil); dsc != "" {
		filters = append(filters, fmt.Sprintf("(description = \"%s\")", dsc))
	}
	if cb := getQualString(d, "created_by", nil); cb != "" {
		filters = append(filters, fmt.Sprintf("(createdBy = \"%s\")", cb))
	}
	if done, ok := d.KeyColumnQuals["done"]; ok {
		filters = append(filters, fmt.Sprintf("(done = %v)", done.GetBoolValue()))
	}

	pageToken := ""
	pageSize := int64(1000)
	for {
		operations, nextPageToken, err := client.ListVPCOperations(ctx, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, op := range operations {
			if len(filters) > 0 {
				if !vpcOperationMatchesFilters(op, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, op)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCOperation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var opID string
	if h != nil && h.Item != nil {
		if op, ok := h.Item.(*VPCOperation); ok {
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
	client := NewVPCClient(tok, 30, cfg)
	op, err := client.GetVPCOperation(ctx, VPCOperationID(opID))
	if err != nil {
		return nil, err
	}
	return op, nil
}

// Manual filter for VPC Operation, since API does not support all filters
func vpcOperationMatchesFilters(op *VPCOperation, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, op.Id) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, op.Description) {
			return false
		}
		if strings.HasPrefix(f, "(createdBy = ") && !strings.Contains(f, op.CreatedBy) {
			return false
		}
		if strings.HasPrefix(f, "(done = ") && !strings.Contains(f, fmt.Sprintf("%v", op.Done)) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only date (YYYY-MM-DD)
func createdAtVPCOperationDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	op, ok := d.HydrateItem.(*VPCOperation)
	if !ok || op.CreatedAt == "" {
		return nil, nil
	}
	if len(op.CreatedAt) < 10 {
		return op.CreatedAt, nil
	}
	return op.CreatedAt[:10], nil
}
