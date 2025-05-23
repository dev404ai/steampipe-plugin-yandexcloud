package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexComputeZone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_compute_zone",
		Description: "Yandex Cloud Compute zones.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"zone_id", "region_id", "status", "name"}),
			Hydrate:    listYandexComputeZones,
		},
		Columns: []*plugin.Column{
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "Zone ID."},
			{Name: "region_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("RegionId"), Description: "Region ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Zone name."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Zone status."},
		},
	}
}

func listYandexComputeZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewComputeClient(tok, 30, cfg)

	var filters []string
	if z := getQualString(d, "zone_id", nil); z != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", z))
	}
	if r := getQualString(d, "region_id", nil); r != "" {
		filters = append(filters, fmt.Sprintf("(regionId = \"%s\")", r))
	}
	if st := getQualString(d, "status", nil); st != "" {
		filters = append(filters, fmt.Sprintf("(status = \"%s\")", strings.ToUpper(st)))
	}
	if n := getQualString(d, "name", nil); n != "" {
		filters = append(filters, fmt.Sprintf("(name = \"%s\")", n))
	}
	// The API does not support filtering, so we filter manually

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
		zones, nextPageToken, err := client.ListZones(ctx, pageToken, pageSize, timeoutSec, retryCount)
		if err != nil {
			return nil, err
		}
		for _, zone := range zones {
			if len(filters) > 0 {
				if !zoneMatchesFilters(zone, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, zone)
		}
		if nextPageToken == "" || nextPageToken == PageToken("") {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

// Manual filtering of zones, since the API does not support filters
func zoneMatchesFilters(zone *Zone, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, zone.Id) {
			return false
		}
		if strings.HasPrefix(f, "(regionId = ") && !strings.Contains(f, zone.RegionId) {
			return false
		}
		if strings.HasPrefix(f, "(status = ") && !strings.Contains(f, zone.Status) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, zone.Name) {
			return false
		}
	}
	return true
}
