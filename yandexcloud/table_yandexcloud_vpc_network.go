package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCNetwork(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_network",
		Description: "Yandex Cloud VPC networks.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "network_id", "name", "description"}),
			Hydrate:    listYandexVPCNetworks,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("network_id"),
			Hydrate:    getYandexVPCNetwork,
		},
		Columns: []*plugin.Column{
			{Name: "network_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "VPC network ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the network."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Network name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Network description."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCNetworkDateTransform), Description: "Network creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
		},
	}
}

func listYandexVPCNetworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if id := getQualString(d, "network_id", nil); id != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", id))
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
		networks, nextPageToken, err := client.ListVPCNetworks(ctx, folderID, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, net := range networks {
			if len(filters) > 0 {
				if !vpcNetworkMatchesFilters(net, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, net)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCNetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var netID string
	if h != nil && h.Item != nil {
		if net, ok := h.Item.(*VPCNetwork); ok {
			netID = net.Id
		}
	}
	if netID == "" {
		if v, ok := d.KeyColumnQuals["network_id"]; ok {
			netID = v.GetStringValue()
		}
	}
	if netID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)
	net, err := client.GetVPCNetwork(ctx, VPCNetworkID(netID))
	if err != nil {
		return nil, err
	}
	return net, nil
}

// Manual filtering, since the API does not support filters except folderId
func vpcNetworkMatchesFilters(net *VPCNetwork, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, net.Id) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, net.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, net.Description) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtVPCNetworkDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	net, ok := d.HydrateItem.(*VPCNetwork)
	if !ok || net.CreatedAt == "" {
		return nil, nil
	}
	if len(net.CreatedAt) < 10 {
		return net.CreatedAt, nil
	}
	return net.CreatedAt[:10], nil
}
