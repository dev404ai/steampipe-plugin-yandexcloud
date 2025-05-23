package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_image",
		Description: "Yandex Cloud Compute disk images.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "status", "name", "family"}),
			Hydrate:    listYandexComputeImages,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("image_id"),
			Hydrate:    getYandexComputeImage,
		},
		Columns: []*plugin.Column{
			{Name: "image_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Image ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Image name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Image description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the image."},
			{Name: "family", Type: proto.ColumnType_STRING, Transform: transform.FromField("Family"), Description: "Image family."},
			{Name: "product_ids", Type: proto.ColumnType_JSON, Transform: transform.FromField("ProductIds"), Description: "Product IDs associated with the image."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtImageDateTransform), Description: "Image creation date (YYYY-MM-DD)."},
			{Name: "min_disk_size", Type: proto.ColumnType_STRING, Transform: transform.FromField("MinDiskSize"), Description: "Minimum disk size required (bytes)."},
			{Name: "size", Type: proto.ColumnType_STRING, Transform: transform.FromField("Size"), Description: "Image size (bytes)."},
			{Name: "os_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("OsType"), Description: "Operating system type."},
			{Name: "os_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("OsVersion"), Description: "Operating system version."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if f := getQualString(d, "family", nil); f != "" {
		filters = append(filters, fmt.Sprintf("(family = \"%s\")", f))
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (images)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (images)", pageToken)
		images, nextPageToken, err := client.ListImages(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, img := range images {
			d.StreamListItem(ctx, img)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var imageID string
	if h != nil && h.Item != nil {
		if img, ok := h.Item.(*Image); ok {
			imageID = img.Id
		}
	}
	if imageID == "" {
		if v, ok := d.KeyColumnQuals["image_id"]; ok {
			imageID = v.GetStringValue()
		}
	}
	if imageID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	img, err := client.GetImage(ctx, ImageID(imageID), 30, 3)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtImageDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	img, ok := d.HydrateItem.(*Image)
	if !ok || img.CreatedAt == "" {
		return nil, nil
	}
	if len(img.CreatedAt) < 10 {
		return img.CreatedAt, nil
	}
	return img.CreatedAt[:10], nil
}
