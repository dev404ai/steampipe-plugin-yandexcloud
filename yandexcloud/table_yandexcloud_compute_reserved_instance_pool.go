package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeReservedInstancePool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_reserved_instance_pool",
		Description: "Yandex Cloud Compute reserved instance pools.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status", "name", "type"}),
			Hydrate:    listYandexComputeReservedInstancePools,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("reserved_instance_pool_id"),
			Hydrate:    getYandexComputeReservedInstancePool,
		},
		Columns: []*plugin.Column{
			{Name: "reserved_instance_pool_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Reserved instance pool ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Reserved instance pool name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Reserved instance pool description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the reserved instance pool."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type"), Description: "Reserved instance pool type."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtReservedInstancePoolDateTransform), Description: "Reserved instance pool creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeReservedInstancePools(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)

	var folderIDStr *string
	if cfg.FolderID != nil {
		str := string(*cfg.FolderID)
		folderIDStr = &str
	}
	folderID := FolderID(getQualString(d, "folder_id", folderIDStr))
	if folderID == "" {
		return nil, fmt.Errorf("folder_id must be provided")
	}

	var filters []string
	if z := getQualString(d, "zone", nil); z != "" {
		filters = append(filters, fmt.Sprintf("(zoneId = \"%s\")", z))
	}
	if st := getQualString(d, "status", nil); st != "" {
		filters = append(filters, fmt.Sprintf("(status = \"%s\")", strings.ToUpper(st)))
	}
	if n := getQualString(d, "name", nil); n != "" {
		filters = append(filters, fmt.Sprintf("(name = \"%s\")", n))
	}
	if t := getQualString(d, "type", nil); t != "" {
		filters = append(filters, fmt.Sprintf("(type = \"%s\")", t))
	}
	filterStr := strings.Join(filters, " AND ")

	filter := Filter(filterStr)
	pageToken := PageToken("")
	pageSize := PageSize(1000)
	timeoutSec := TimeoutSec(30)
	if cfg.Timeout != nil && *cfg.Timeout > 0 {
		timeoutSec = TimeoutSec(*cfg.Timeout)
	}
	retryCount := RetryCount(3)
	if cfg.Retry != nil && *cfg.Retry > 0 {
		retryCount = RetryCount(*cfg.Retry)
	}
	LogInfo(ctx, "DEBUG: folderID used: '%s' (reserved instance pools)", string(folderID))
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (reserved instance pools)", string(pageToken))
		pools, nextPageToken, err := client.ListReservedInstancePools(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, pool := range pools {
			d.StreamListItem(ctx, pool)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeReservedInstancePool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var poolID string
	if h != nil && h.Item != nil {
		if pool, ok := h.Item.(*ReservedInstancePool); ok {
			poolID = pool.Id
		}
	}
	if poolID == "" {
		if v, ok := d.KeyColumnQuals["reserved_instance_pool_id"]; ok {
			poolID = v.GetStringValue()
		}
	}
	if poolID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	pool, err := client.GetReservedInstancePool(ctx, ReservedInstancePoolID(poolID), 30, 3)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtReservedInstancePoolDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	pool, ok := d.HydrateItem.(*ReservedInstancePool)
	if !ok || pool.CreatedAt == "" {
		return nil, nil
	}
	if len(pool.CreatedAt) < 10 {
		return pool.CreatedAt, nil
	}
	return pool.CreatedAt[:10], nil
}
