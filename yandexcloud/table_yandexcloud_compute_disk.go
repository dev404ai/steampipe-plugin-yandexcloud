package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeDisk(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_disk",
		Description: "Yandex Cloud Compute disks.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status", "name", "type_id"}),
			Hydrate:    listYandexComputeDisks,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("disk_id"),
			Hydrate:    getYandexComputeDisk,
		},
		Columns: []*plugin.Column{
			{Name: "disk_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Disk ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Disk name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Disk description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the disk."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "type_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("TypeId"), Description: "Disk type ID."},
			{Name: "size", Type: proto.ColumnType_STRING, Transform: transform.FromField("Size"), Description: "Disk size (bytes)."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtDiskDateTransform), Description: "Disk creation date (YYYY-MM-DD)."},
			{Name: "source_image_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("SourceImageId"), Description: "Source image ID."},
			{Name: "source_snapshot_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("SourceSnapshotId"), Description: "Source snapshot ID."},
			{Name: "block_size", Type: proto.ColumnType_STRING, Transform: transform.FromField("BlockSize"), Description: "Block size (bytes)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeDisks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		LogError(ctx, "Failed to get token: %v", err)
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
	if t := getQualString(d, "type_id", nil); t != "" {
		filters = append(filters, fmt.Sprintf("(typeId = \"%s\")", t))
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (disks)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (disks)", pageToken)
		disks, nextPageToken, err := client.ListDisks(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, disk := range disks {
			d.StreamListItem(ctx, disk)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeDisk(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var diskID string
	if h != nil && h.Item != nil {
		if disk, ok := h.Item.(*Disk); ok {
			diskID = disk.Id
		}
	}
	if diskID == "" {
		if v, ok := d.KeyColumnQuals["disk_id"]; ok {
			diskID = v.GetStringValue()
		}
	}
	if diskID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	disk, err := client.GetDisk(ctx, DiskID(diskID), 30, 3)
	if err != nil {
		return nil, err
	}
	return disk, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtDiskDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	disk, ok := d.HydrateItem.(*Disk)
	if !ok || disk.CreatedAt == "" {
		return nil, nil
	}
	if len(disk.CreatedAt) < 10 {
		return disk.CreatedAt, nil
	}
	return disk.CreatedAt[:10], nil
}
