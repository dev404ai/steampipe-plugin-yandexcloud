package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_snapshot",
		Description: "Yandex Cloud Compute disk snapshots.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status", "name"}),
			Hydrate:    listYandexComputeSnapshots,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("snapshot_id"),
			Hydrate:    getYandexComputeSnapshot,
		},
		Columns: []*plugin.Column{
			{Name: "snapshot_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Snapshot ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Snapshot name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Snapshot description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the snapshot."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtSnapshotDateTransform), Description: "Snapshot creation date (YYYY-MM-DD)."},
			{Name: "source_disk_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("SourceDiskId"), Description: "Source disk ID."},
			{Name: "size", Type: proto.ColumnType_STRING, Transform: transform.FromField("Size"), Description: "Snapshot size (bytes)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (snapshots)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (snapshots)", pageToken)
		snapshots, nextPageToken, err := client.ListSnapshots(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, snap := range snapshots {
			d.StreamListItem(ctx, snap)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var snapshotID string
	if h != nil && h.Item != nil {
		if snap, ok := h.Item.(*Snapshot); ok {
			snapshotID = snap.Id
		}
	}
	if snapshotID == "" {
		if v, ok := d.KeyColumnQuals["snapshot_id"]; ok {
			snapshotID = v.GetStringValue()
		}
	}
	if snapshotID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	snap, err := client.GetSnapshot(ctx, SnapshotID(snapshotID), 30, 3)
	if err != nil {
		return nil, err
	}
	return snap, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtSnapshotDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	snap, ok := d.HydrateItem.(*Snapshot)
	if !ok || snap.CreatedAt == "" {
		return nil, nil
	}
	if len(snap.CreatedAt) < 10 {
		return snap.CreatedAt, nil
	}
	return snap.CreatedAt[:10], nil
}
