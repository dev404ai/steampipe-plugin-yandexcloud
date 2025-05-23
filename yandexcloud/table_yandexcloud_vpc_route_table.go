package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCRouteTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_route_table",
		Description: "Yandex Cloud VPC route tables.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "route_table_id", "network_id", "name", "description"}),
			Hydrate:    listYandexVPCRouteTables,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("route_table_id"),
			Hydrate:    getYandexVPCRouteTable,
		},
		Columns: []*plugin.Column{
			{Name: "route_table_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "VPC route table ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the route table."},
			{Name: "network_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NetworkId"), Description: "Network ID to which the route table belongs."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Route table name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Route table description."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCRouteTableDateTransform), Description: "Route table creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
			{Name: "static_routes", Type: proto.ColumnType_JSON, Transform: transform.FromField("StaticRoutes"), Description: "List of static routes."},
		},
	}
}

func listYandexVPCRouteTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)

	var folderIDStr *string
	if cfg.FolderID != nil {
		str := string(*cfg.FolderID)
		folderIDStr = &str
	}
	folderID := getQualString(d, "folder_id", folderIDStr)
	if folderID == "" {
		return nil, fmt.Errorf("folder_id must be provided")
	}

	var filters []string
	if id := getQualString(d, "route_table_id", nil); id != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", id))
	}
	if netid := getQualString(d, "network_id", nil); netid != "" {
		filters = append(filters, fmt.Sprintf("(networkId = \"%s\")", netid))
	}
	if n := getQualString(d, "name", nil); n != "" {
		filters = append(filters, fmt.Sprintf("(name = \"%s\")", n))
	}
	if dsc := getQualString(d, "description", nil); dsc != "" {
		filters = append(filters, fmt.Sprintf("(description = \"%s\")", dsc))
	}

	pageToken := ""
	pageSize := int64(1000)
	for {
		routeTables, nextPageToken, err := client.ListVPCRouteTables(ctx, folderID, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, rt := range routeTables {
			if len(filters) > 0 {
				if !vpcRouteTableMatchesFilters(rt, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, rt)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCRouteTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var rtID string
	if h != nil && h.Item != nil {
		if rt, ok := h.Item.(*VPCRouteTable); ok {
			rtID = rt.Id
		}
	}
	if rtID == "" {
		if v, ok := d.KeyColumnQuals["route_table_id"]; ok {
			rtID = v.GetStringValue()
		}
	}
	if rtID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)
	rt, err := client.GetVPCRouteTable(ctx, VPCRouteTableID(rtID))
	if err != nil {
		return nil, err
	}
	return rt, nil
}

// Manual filtering, since the API does not support filters except folderId
func vpcRouteTableMatchesFilters(rt *VPCRouteTable, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, rt.Id) {
			return false
		}
		if strings.HasPrefix(f, "(networkId = ") && !strings.Contains(f, rt.NetworkId) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, rt.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, rt.Description) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtVPCRouteTableDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	rt, ok := d.HydrateItem.(*VPCRouteTable)
	if !ok || rt.CreatedAt == "" {
		return nil, nil
	}
	if len(rt.CreatedAt) < 10 {
		return rt.CreatedAt, nil
	}
	return rt.CreatedAt[:10], nil
}
