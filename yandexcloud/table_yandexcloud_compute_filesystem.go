package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeFilesystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_filesystem",
		Description: "Yandex Cloud Compute filesystems.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status", "name", "type_id"}),
			Hydrate:    listYandexComputeFilesystems,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("filesystem_id"),
			Hydrate:    getYandexComputeFilesystem,
		},
		Columns: []*plugin.Column{
			{Name: "filesystem_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Filesystem ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Filesystem name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Filesystem description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the filesystem."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "type_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("TypeId"), Description: "Filesystem type ID."},
			{Name: "size", Type: proto.ColumnType_STRING, Transform: transform.FromField("Size"), Description: "Filesystem size (bytes)."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtFilesystemDateTransform), Description: "Filesystem creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeFilesystems(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (filesystems)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (filesystems)", pageToken)
		filesystems, nextPageToken, err := client.ListFilesystems(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, fs := range filesystems {
			d.StreamListItem(ctx, fs)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeFilesystem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var filesystemID string
	if h != nil && h.Item != nil {
		if fs, ok := h.Item.(*Filesystem); ok {
			filesystemID = fs.Id
		}
	}
	if filesystemID == "" {
		if v, ok := d.KeyColumnQuals["filesystem_id"]; ok {
			filesystemID = v.GetStringValue()
		}
	}
	if filesystemID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	fs, err := client.GetFilesystem(ctx, FilesystemID(filesystemID), 30, 3)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtFilesystemDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	fs, ok := d.HydrateItem.(*Filesystem)
	if !ok || fs.CreatedAt == "" {
		return nil, nil
	}
	if len(fs.CreatedAt) < 10 {
		return fs.CreatedAt, nil
	}
	return fs.CreatedAt[:10], nil
}
