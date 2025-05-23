package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeGPUCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_gpu_cluster",
		Description: "Yandex Cloud Compute GPU clusters.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status", "name", "type"}),
			Hydrate:    listYandexComputeGPUClusters,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("gpu_cluster_id"),
			Hydrate:    getYandexComputeGPUCluster,
		},
		Columns: []*plugin.Column{
			{Name: "gpu_cluster_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "GPU cluster ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "GPU cluster name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "GPU cluster description."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the GPU cluster."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type"), Description: "GPU cluster type."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtGPUClusterDateTransform), Description: "GPU cluster creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexComputeGPUClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	LogInfo(ctx, "DEBUG: folderID used: '%s' (gpu clusters)", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s' (gpu clusters)", pageToken)
		clusters, nextPageToken, err := client.ListGPUClusters(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, cluster := range clusters {
			d.StreamListItem(ctx, cluster)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexComputeGPUCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterID string
	if h != nil && h.Item != nil {
		if cluster, ok := h.Item.(*GPUCluster); ok {
			clusterID = cluster.Id
		}
	}
	if clusterID == "" {
		if v, ok := d.KeyColumnQuals["gpu_cluster_id"]; ok {
			clusterID = v.GetStringValue()
		}
	}
	if clusterID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	cluster, err := client.GetGPUCluster(ctx, GPUClusterID(clusterID), 30, 3)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtGPUClusterDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	cluster, ok := d.HydrateItem.(*GPUCluster)
	if !ok || cluster.CreatedAt == "" {
		return nil, nil
	}
	if len(cluster.CreatedAt) < 10 {
		return cluster.CreatedAt, nil
	}
	return cluster.CreatedAt[:10], nil
}
