package yandexcloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

// Plugin returns the top-level plugin definition used by Steampipe.
func Plugin() *plugin.Plugin {
	ctx := context.Background()
	return &plugin.Plugin{
		Name:                   "yandexcloud",
		ConnectionConfigSchema: connectionConfig(),
		TableMap: map[string]*plugin.Table{
			"yandexcloud_compute_instance":               tableYandexComputeInstance(ctx),
			"yandexcloud_billing_resource_usage":         tableYandexBillingResourceUsage(ctx),
			"yandexcloud_compute_snapshot":               tableYandexComputeSnapshot(ctx),
			"yandexcloud_compute_image":                  tableYandexComputeImage(ctx),
			"yandexcloud_compute_disk":                   tableYandexComputeDisk(ctx),
			"yandexcloud_compute_filesystem":             tableYandexComputeFilesystem(ctx),
			"yandexcloud_compute_placement_group":        tableYandexComputePlacementGroup(ctx),
			"yandexcloud_compute_host_group":             tableYandexComputeHostGroup(ctx),
			"yandexcloud_compute_gpu_cluster":            tableYandexComputeGPUCluster(ctx),
			"yandexcloud_compute_disk_placement_group":   tableYandexComputeDiskPlacementGroup(ctx),
			"yandexcloud_compute_snapshot_schedule":      tableYandexComputeSnapshotSchedule(ctx),
			"yandexcloud_compute_reserved_instance_pool": tableYandexComputeReservedInstancePool(ctx),
			"yandexcloud_compute_zone":                   tableYandexComputeZone(ctx),
			"yandexcloud_compute_disk_type":              tableYandexComputeDiskType(ctx),
			"yandexcloud_compute_host_type":              tableYandexComputeHostType(ctx),
			"yandexcloud_compute_operation":              tableYandexComputeOperation(ctx),
			"yandexcloud_vpc_network":                    tableYandexVPCNetwork(ctx),
			"yandexcloud_vpc_subnet":                     tableYandexVPCSubnet(ctx),
			"yandexcloud_vpc_route_table":                tableYandexVPCRouteTable(ctx),
			"yandexcloud_vpc_security_group":             tableYandexVPCSecurityGroup(ctx),
			"yandexcloud_vpc_address":                    tableYandexVPCAddress(ctx),
			"yandexcloud_vpc_gateway":                    tableYandexVPCGateway(ctx),
			"yandexcloud_vpc_operation":                  tableYandexVPCOperation(ctx),
			"yandexcloud_billing_account":                tableYandexBillingAccount(ctx),
			"yandexcloud_billing_sku":                    tableYandexBillingSku(ctx),
			"yandexcloud_billing_budget":                 tableYandexBillingBudget(ctx),
		},
	}
}
