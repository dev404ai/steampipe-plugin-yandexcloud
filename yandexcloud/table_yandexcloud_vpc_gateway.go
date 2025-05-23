package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_gateway",
		Description: "Yandex Cloud VPC gateways.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "gateway_id", "name", "description"}),
			Hydrate:    listYandexVPCGateways,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("gateway_id"),
			Hydrate:    getYandexVPCGateway,
		},
		Columns: []*plugin.Column{
			{Name: "gateway_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "VPC gateway ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the gateway."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCGatewayDateTransform), Description: "Gateway creation date (YYYY-MM-DD)."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Gateway name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Gateway description."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
			{Name: "shared_egress_gateway", Type: proto.ColumnType_JSON, Transform: transform.FromField("SharedEgressGateway"), Description: "Shared egress gateway specification."},
		},
	}
}

func listYandexVPCGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if id := getQualString(d, "gateway_id", nil); id != "" {
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
		gateways, nextPageToken, err := client.ListVPCGateways(ctx, folderID, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, gw := range gateways {
			if len(filters) > 0 {
				if !vpcGatewayMatchesFilters(gw, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, gw)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var gwID string
	if h != nil && h.Item != nil {
		if gw, ok := h.Item.(*VPCGateway); ok {
			gwID = gw.Id
		}
	}
	if gwID == "" {
		if v, ok := d.KeyColumnQuals["gateway_id"]; ok {
			gwID = v.GetStringValue()
		}
	}
	if gwID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)
	gw, err := client.GetVPCGateway(ctx, VPCGatewayID(gwID))
	if err != nil {
		return nil, err
	}
	return gw, nil
}

// Manual filter for VPC Gateway, since API does not support all filters
func vpcGatewayMatchesFilters(gw *VPCGateway, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, gw.Id) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, gw.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, gw.Description) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only date (YYYY-MM-DD)
func createdAtVPCGatewayDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	gw, ok := d.HydrateItem.(*VPCGateway)
	if !ok || gw.CreatedAt == "" {
		return nil, nil
	}
	if len(gw.CreatedAt) < 10 {
		return gw.CreatedAt, nil
	}
	return gw.CreatedAt[:10], nil
}
