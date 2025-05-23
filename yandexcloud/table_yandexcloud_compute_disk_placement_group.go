package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeDiskPlacementGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_disk_placement_group",
		Description: "Yandex Cloud Compute disk placement groups.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status", "name", "type"}),
			Hydrate:    listYandexComputeDiskPlacementGroups,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("disk_placement_group_id"),
			Hydrate:    getYandexComputeDiskPlacementGroup,
		},
		Columns: []*plugin.Column{
			{Name: "disk_placement_group_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Disk placement group ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Disk placement group name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Disk placement group description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the disk placement group."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type"), Description: "Disk placement group type."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtDiskPlacementGroupDateTransform), Description: "Disk placement group creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeDiskPlacementGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (disk placement groups)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (disk placement groups)", pageToken)
		groups, nextPageToken, err := client.ListDiskPlacementGroups(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, group := range groups {
			d.StreamListItem(ctx, group)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeDiskPlacementGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var groupID string
	if h != nil && h.Item != nil {
		if group, ok := h.Item.(*DiskPlacementGroup); ok {
			groupID = group.Id
		}
	}
	if groupID == "" {
		if v, ok := d.KeyColumnQuals["disk_placement_group_id"]; ok {
			groupID = v.GetStringValue()
		}
	}
	if groupID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	group, err := client.GetDiskPlacementGroup(ctx, DiskPlacementGroupID(groupID), 30, 3)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtDiskPlacementGroupDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	group, ok := d.HydrateItem.(*DiskPlacementGroup)
	if !ok || group.CreatedAt == "" {
		return nil, nil
	}
	if len(group.CreatedAt) < 10 {
		return group.CreatedAt, nil
	}
	return group.CreatedAt[:10], nil
}
