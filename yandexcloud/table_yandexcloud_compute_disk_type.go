package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeDiskType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_disk_type",
		Description: "Yandex Cloud Compute disk types.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"zone_id", "disk_type_id", "name", "description"}),
			Hydrate:    listYandexComputeDiskTypes,
		},
		Columns: []*plugin.Column{
			{Name: "disk_type_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Disk type ID."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Zone ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Disk type name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Disk type description."},
		},
	}
}

func listYandexComputeDiskTypes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)

	zoneID := getQualString(d, "zone_id", nil)
	var filters []string
	if id := getQualString(d, "disk_type_id", nil); id != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", id))
	}
	if n := getQualString(d, "name", nil); n != "" {
		filters = append(filters, fmt.Sprintf("(name = \"%s\")", n))
	}
	if dsc := getQualString(d, "description", nil); dsc != "" {
		filters = append(filters, fmt.Sprintf("(description = \"%s\")", dsc))
	}
	// The API supports filtering only by zoneId, other filters are manual

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
	for {
		diskTypes, nextPageToken, err := client.ListDiskTypes(ctx, zoneID, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, dt := range diskTypes {
			if len(filters) > 0 {
				if !diskTypeMatchesFilters(dt, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, dt)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

// Manual filtering of disk types, since the API does not support filters except zoneId
func diskTypeMatchesFilters(dt *DiskType, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, dt.Id) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, dt.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, dt.Description) {
			return false
		}
	}
	return true
}
