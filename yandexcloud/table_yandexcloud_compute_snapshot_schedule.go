package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeSnapshotSchedule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_snapshot_schedule",
		Description: "Yandex Cloud Compute snapshot schedules.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "status", "name", "description"}),
			Hydrate:    listYandexComputeSnapshotSchedules,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("snapshot_schedule_id"),
			Hydrate:    getYandexComputeSnapshotSchedule,
		},
		Columns: []*plugin.Column{
			{Name: "snapshot_schedule_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Snapshot schedule ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Snapshot schedule name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Snapshot schedule description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the snapshot schedule."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtSnapshotScheduleDateTransform), Description: "Snapshot schedule creation date (YYYY-MM-DD)."},
			{Name: "schedule_policy", Type: proto.ColumnType_JSON, Transform: transform.FromField("SchedulePolicy"), Description: "Schedule policy details."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeSnapshotSchedules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if st := getQualString(d, "status", nil); st != "" {
		filters = append(filters, fmt.Sprintf("(status = \"%s\")", strings.ToUpper(st)))
	}
	if n := getQualString(d, "name", nil); n != "" {
		filters = append(filters, fmt.Sprintf("(name = \"%s\")", n))
	}
	if dsc := getQualString(d, "description", nil); dsc != "" {
		filters = append(filters, fmt.Sprintf("(description = \"%s\")", dsc))
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (snapshot schedules)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (snapshot schedules)", pageToken)
		schedules, nextPageToken, err := client.ListSnapshotSchedules(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, schedule := range schedules {
			d.StreamListItem(ctx, schedule)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeSnapshotSchedule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var scheduleID string
	if h != nil && h.Item != nil {
		if schedule, ok := h.Item.(*SnapshotSchedule); ok {
			scheduleID = schedule.Id
		}
	}
	if scheduleID == "" {
		if v, ok := d.KeyColumnQuals["snapshot_schedule_id"]; ok {
			scheduleID = v.GetStringValue()
		}
	}
	if scheduleID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	schedule, err := client.GetSnapshotSchedule(ctx, SnapshotScheduleID(scheduleID), 30, 3)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtSnapshotScheduleDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	schedule, ok := d.HydrateItem.(*SnapshotSchedule)
	if !ok || schedule.CreatedAt == "" {
		return nil, nil
	}
	if len(schedule.CreatedAt) < 10 {
		return schedule.CreatedAt, nil
	}
	return schedule.CreatedAt[:10], nil
}
