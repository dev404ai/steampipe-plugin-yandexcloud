package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCSubnet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_subnet",
		Description: "Yandex Cloud VPC subnets.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "subnet_id", "network_id", "zone_id", "name", "description"}),
			Hydrate:    listYandexVPCSubnets,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("subnet_id"),
			Hydrate:    getYandexVPCSubnet,
		},
		Columns: []*plugin.Column{
			{Name: "subnet_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "VPC subnet ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the subnet."},
			{Name: "network_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NetworkId"), Description: "Network ID to which the subnet belongs."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneId"), Description: "Zone ID where the subnet is located."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Subnet name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Subnet description."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCSubnetDateTransform), Description: "Subnet creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
			{Name: "cidr_blocks", Type: proto.ColumnType_JSON, Transform: transform.FromField("CidrBlocks"), Description: "List of IPv4 CIDR blocks assigned to the subnet."},
		},
	}
}

func listYandexVPCSubnets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if id := getQualString(d, "subnet_id", nil); id != "" {
		filters = append(filters, fmt.Sprintf("(id = \"%s\")", id))
	}
	if netid := getQualString(d, "network_id", nil); netid != "" {
		filters = append(filters, fmt.Sprintf("(networkId = \"%s\")", netid))
	}
	if zone := getQualString(d, "zone_id", nil); zone != "" {
		filters = append(filters, fmt.Sprintf("(zoneId = \"%s\")", zone))
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
		subnets, nextPageToken, err := client.ListVPCSubnets(ctx, folderID, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, subnet := range subnets {
			if len(filters) > 0 {
				if !vpcSubnetMatchesFilters(subnet, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, subnet)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCSubnet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var subnetID string
	if h != nil && h.Item != nil {
		if subnet, ok := h.Item.(*VPCSubnet); ok {
			subnetID = subnet.Id
		}
	}
	if subnetID == "" {
		if v, ok := d.KeyColumnQuals["subnet_id"]; ok {
			subnetID = v.GetStringValue()
		}
	}
	if subnetID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)
	subnet, err := client.GetVPCSubnet(ctx, VPCSubnetID(subnetID))
	if err != nil {
		return nil, err
	}
	return subnet, nil
}

// Manual filtering, since the API does not support filters except folderId
func vpcSubnetMatchesFilters(subnet *VPCSubnet, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, subnet.Id) {
			return false
		}
		if strings.HasPrefix(f, "(networkId = ") && !strings.Contains(f, subnet.NetworkId) {
			return false
		}
		if strings.HasPrefix(f, "(zoneId = ") && !strings.Contains(f, subnet.ZoneId) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, subnet.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, subnet.Description) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtVPCSubnetDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	subnet, ok := d.HydrateItem.(*VPCSubnet)
	if !ok || subnet.CreatedAt == "" {
		return nil, nil
	}
	if len(subnet.CreatedAt) < 10 {
		return subnet.CreatedAt, nil
	}
	return subnet.CreatedAt[:10], nil
}
