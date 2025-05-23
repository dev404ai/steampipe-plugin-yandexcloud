package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_instance",
		Description: "Yandex Cloud Compute virtual machine instances.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "zone", "status"}),
			Hydrate:    listYandexComputeInstances,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("instance_id"),
			Hydrate:    getYandexComputeInstance,
		},
		Columns: []*plugin.Column{
			{Name: "instance_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Instance ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Instance name."},
			{Name: "zone", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Availability zone."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the instance."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Instance description."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
			{Name: "platform_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("PlatformId"), Description: "Hardware platform configuration ID."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtDateTransform), Description: "Instance creation date (YYYY-MM-DD)."},
			{Name: "resources", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resources"), Description: "Computing resources (CPU, RAM, GPU, core_fraction)."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Hydrate: getYandexComputeInstance, Transform: transform.FromField("Metadata"), Description: "Instance metadata (e.g., ssh-keys)."},
			{Name: "metadata_options", Type: proto.ColumnType_JSON, Hydrate: getYandexComputeInstance, Transform: transform.FromField("MetadataOptions"), Description: "Metadata access options."},
			{Name: "boot_disk", Type: proto.ColumnType_JSON, Transform: transform.FromField("BootDisk"), Description: "Boot disk information."},
			{Name: "secondary_disks", Type: proto.ColumnType_JSON, Transform: transform.FromField("SecondaryDisks"), Description: "Secondary disks attached to the instance."},
			{Name: "network_interfaces", Type: proto.ColumnType_JSON, Transform: transform.FromField("NetworkInterfaces"), Description: "Network interfaces attached to the instance."},
			{Name: "fqdn", Type: proto.ColumnType_STRING, Hydrate: getYandexComputeInstance, Transform: transform.FromField("FQDN"), Description: "Fully qualified domain name."},
			{Name: "service_account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ServiceAccountId"), Description: "Service account ID attached to the instance."},
			{Name: "hostname", Type: proto.ColumnType_STRING, Transform: transform.FromField("Hostname"), Description: "Instance hostname."},
			{Name: "deletion_protection", Type: proto.ColumnType_BOOL, Transform: transform.FromField("DeletionProtection"), Description: "Deletion protection flag."},
		},
	}
}

func listYandexComputeInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

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
	LogInfo(ctx, "DEBUG: folderID used: '%s'", folderID)
	for {
		LogInfo(ctx, "DEBUG: pageToken: '%s'", pageToken)
		instances, nextPageToken, err := client.ListInstances(ctx, folderID, filter, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, inst := range instances {
			d.StreamListItem(ctx, inst)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

// Hydrate function to get a single VM by instance_id
func getYandexComputeInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var instanceID string
	if h != nil && h.Item != nil {
		if inst, ok := h.Item.(*Instance); ok {
			instanceID = inst.Id
		}
	}
	if instanceID == "" {
		if v, ok := d.KeyColumnQuals["instance_id"]; ok {
			instanceID = v.GetStringValue()
		}
	}
	if instanceID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)
	inst, err := client.GetInstance(ctx, InstanceID(instanceID), 30, 3)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	inst, ok := d.HydrateItem.(*Instance)
	if !ok || inst.CreatedAt == "" {
		return nil, nil
	}
	if len(inst.CreatedAt) < 10 {
		return inst.CreatedAt, nil
	}
	return inst.CreatedAt[:10], nil
}
