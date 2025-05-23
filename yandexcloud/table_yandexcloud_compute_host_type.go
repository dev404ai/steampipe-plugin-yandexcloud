package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeHostType(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_host_type",
		Description: "Yandex Cloud Compute host types.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"zone_id", "host_type_id", "name", "description"}),
			Hydrate:    listYandexComputeHostTypes,
		},
		Columns: []*plugin.Column{
			{Name: "host_type_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Host type ID."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Zone ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Host type name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Host type description."},
		},
	}
}

func listYandexComputeHostTypes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)

	zoneID := getQualString(d, "zone_id", nil)
	var filters []string
	if id := getQualString(d, "host_type_id", nil); id != "" {
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
		hostTypes, nextPageToken, err := client.ListHostTypes(ctx, zoneID, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, ht := range hostTypes {
			if len(filters) > 0 {
				if !hostTypeMatchesFilters(ht, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, ht)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

// Manual filtering of host types, since the API does not support filters except zoneId
func hostTypeMatchesFilters(ht *HostType, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, ht.Id) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, ht.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, ht.Description) {
			return false
		}
	}
	return true
}
