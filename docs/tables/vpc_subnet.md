---
title: Table: yandexcloud_vpc_subnet
summary: Query information about Yandex Cloud VPC subnets.
---

# Table: yandexcloud_vpc_subnet

The `yandexcloud_vpc_subnet` table allows you to query information about subnets in Yandex Cloud VPC.

## Examples

### List all subnets
```sql
select subnet_id, name, network_id, zone_id, created_at from yandexcloud_vpc_subnet;
```

### Find subnets by network
```sql
select subnet_id, name from yandexcloud_vpc_subnet where network_id = 'network-123';
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| subnet_id   | text   | VPC subnet ID.                              |
| folder_id   | text   | Folder ID containing the subnet.            |
| network_id  | text   | Network ID to which the subnet belongs.     |
| zone_id     | text   | Zone ID where the subnet is located.        |
| name        | text   | Subnet name.                                |
| description | text   | Subnet description.                         |
| created_at  | text   | Subnet creation date (YYYY-MM-DD).          |
| labels      | jsonb  | Resource labels as key:value pairs.         | 