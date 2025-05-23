package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCSecurityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_security_group",
		Description: "Yandex Cloud VPC security groups.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "security_group_id", "network_id", "name", "description"}),
			Hydrate:    listYandexVPCSecurityGroups,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("security_group_id"),
			Hydrate:    getYandexVPCSecurityGroup,
		},
		Columns: []*plugin.Column{
			{Name: "security_group_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "VPC security group ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the security group."},
			{Name: "network_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NetworkId"), Description: "Network ID to which the security group belongs."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Security group name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Security group description."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCSecurityGroupDateTransform), Description: "Security group creation date (YYYY-MM-DD)."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
			{Name: "rules", Type: proto.ColumnType_JSON, Transform: transform.FromField("Rules"), Description: "All rules (deprecated, use ingress/egress)."},
			{Name: "ingress_rules", Type: proto.ColumnType_JSON, Transform: transform.FromField("IngressRules"), Description: "Ingress rules."},
			{Name: "egress_rules", Type: proto.ColumnType_JSON, Transform: transform.FromField("EgressRules"), Description: "Egress rules."},
		},
	}
}

func listYandexVPCSecurityGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if id := getQualString(d, "security_group_id", nil); id != "" {
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
		groups, nextPageToken, err := client.ListVPCSecurityGroups(ctx, folderID, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, group := range groups {
			if len(filters) > 0 {
				if !vpcSecurityGroupMatchesFilters(group, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, group)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var groupID string
	if h != nil && h.Item != nil {
		if group, ok := h.Item.(*VPCSecurityGroup); ok {
			groupID = group.Id
		}
	}
	if groupID == "" {
		if v, ok := d.KeyColumnQuals["security_group_id"]; ok {
			groupID = v.GetStringValue()
		}
	}
	if groupID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)
	group, err := client.GetVPCSecurityGroup(ctx, VPCSecurityGroupID(groupID))
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Manual filtering, since the API does not support filters except folderId
func vpcSecurityGroupMatchesFilters(group *VPCSecurityGroup, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, group.Id) {
			return false
		}
		if strings.HasPrefix(f, "(networkId = ") && !strings.Contains(f, group.NetworkId) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, group.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, group.Description) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only the date (YYYY-MM-DD)
func createdAtVPCSecurityGroupDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	group, ok := d.HydrateItem.(*VPCSecurityGroup)
	if !ok || group.CreatedAt == "" {
		return nil, nil
	}
	if len(group.CreatedAt) < 10 {
		return group.CreatedAt, nil
	}
	return group.CreatedAt[:10], nil
}
