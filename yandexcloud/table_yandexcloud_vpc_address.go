package yandexcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableYandexVPCAddress(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yandexcloud_vpc_address",
		Description: "Yandex Cloud VPC addresses.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"folder_id", "address_id", "name", "description", "type", "ip_version", "reserved", "used", "deletion_protection"}),
			Hydrate:    listYandexVPCAddresses,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("address_id"),
			Hydrate:    getYandexVPCAddress,
		},
		Columns: []*plugin.Column{
			{Name: "address_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "VPC address ID."},
			{Name: "folder_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("FolderId"), Description: "Folder ID containing the address."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.From(createdAtVPCAddressDateTransform), Description: "Address creation date (YYYY-MM-DD)."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Address name."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "Address description."},
			{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "Resource labels as key:value pairs."},
			{Name: "external_ipv4_address", Type: proto.ColumnType_JSON, Transform: transform.FromField("ExternalIpv4Address"), Description: "External IPv4 address specification."},
			{Name: "reserved", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Reserved"), Description: "Specifies if address is reserved."},
			{Name: "used", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Used"), Description: "Specifies if address is used."},
			{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type"), Description: "Type of the IP address (INTERNAL/EXTERNAL)."},
			{Name: "ip_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("IpVersion"), Description: "Version of the IP address (IPV4/IPV6)."},
			{Name: "deletion_protection", Type: proto.ColumnType_BOOL, Transform: transform.FromField("DeletionProtection"), Description: "Specifies if address is protected from deletion."},
			{Name: "dns_records", Type: proto.ColumnType_JSON, Transform: transform.FromField("DnsRecords"), Description: "DNS record specifications."},
		},
	}
}

func listYandexVPCAddresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	if id := getQualString(d, "address_id", nil); id != "" {
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
	if t := getQualString(d, "type", nil); t != "" {
		filters = append(filters, fmt.Sprintf("(type = \"%s\")", t))
	}
	if v := getQualString(d, "ip_version", nil); v != "" {
		filters = append(filters, fmt.Sprintf("(ipVersion = \"%s\")", v))
	}
	if r, ok := d.KeyColumnQuals["reserved"]; ok {
		filters = append(filters, fmt.Sprintf("(reserved = %v)", r.GetBoolValue()))
	}
	if u, ok := d.KeyColumnQuals["used"]; ok {
		filters = append(filters, fmt.Sprintf("(used = %v)", u.GetBoolValue()))
	}
	if dp, ok := d.KeyColumnQuals["deletion_protection"]; ok {
		filters = append(filters, fmt.Sprintf("(deletionProtection = %v)", dp.GetBoolValue()))
	}

	pageToken := ""
	pageSize := int64(1000)
	for {
		addresses, nextPageToken, err := client.ListVPCAddresses(ctx, folderID, pageToken, pageSize)
		if err != nil {
			return nil, err
		}
		for _, addr := range addresses {
			if len(filters) > 0 {
				if !vpcAddressMatchesFilters(addr, filters) {
					continue
				}
			}
			d.StreamListItem(ctx, addr)
		}
		if nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return nil, nil
}

func getYandexVPCAddress(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var addrID string
	if h != nil && h.Item != nil {
		if addr, ok := h.Item.(*VPCAddress); ok {
			addrID = addr.Id
		}
	}
	if addrID == "" {
		if v, ok := d.KeyColumnQuals["address_id"]; ok {
			addrID = v.GetStringValue()
		}
	}
	if addrID == "" {
		return nil, nil
	}
	cfg := getConfig(d)
	tok, err := getAuthToken(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := NewVPCClient(tok, 30, cfg)
	addr, err := client.GetVPCAddress(ctx, VPCAddressID(addrID))
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// Manual filter for VPC Address, since API does not support all filters
func vpcAddressMatchesFilters(addr *VPCAddress, filters []string) bool {
	for _, f := range filters {
		if strings.HasPrefix(f, "(id = ") && !strings.Contains(f, addr.Id) {
			return false
		}
		if strings.HasPrefix(f, "(name = ") && !strings.Contains(f, addr.Name) {
			return false
		}
		if strings.HasPrefix(f, "(description = ") && !strings.Contains(f, addr.Description) {
			return false
		}
		if strings.HasPrefix(f, "(type = ") && !strings.Contains(f, addr.Type) {
			return false
		}
		if strings.HasPrefix(f, "(ipVersion = ") && !strings.Contains(f, addr.IpVersion) {
			return false
		}
		if strings.HasPrefix(f, "(reserved = ") && !strings.Contains(f, fmt.Sprintf("%v", addr.Reserved)) {
			return false
		}
		if strings.HasPrefix(f, "(used = ") && !strings.Contains(f, fmt.Sprintf("%v", addr.Used)) {
			return false
		}
		if strings.HasPrefix(f, "(deletionProtection = ") && !strings.Contains(f, fmt.Sprintf("%v", addr.DeletionProtection)) {
			return false
		}
	}
	return true
}

// Transform function for created_at: returns only date (YYYY-MM-DD)
func createdAtVPCAddressDateTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.HydrateItem == nil {
		return nil, nil
	}
	addr, ok := d.HydrateItem.(*VPCAddress)
	if !ok || addr.CreatedAt == "" {
		return nil, nil
	}
	if len(addr.CreatedAt) < 10 {
		return addr.CreatedAt, nil
	}
	return addr.CreatedAt[:10], nil
}
