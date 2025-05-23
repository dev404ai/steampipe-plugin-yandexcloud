---
title: Table: yandexcloud_vpc_network
summary: Query information about Yandex Cloud VPC networks.
---

# Table: yandexcloud_vpc_network

The `yandexcloud_vpc_network` table allows you to query information about VPC networks in Yandex Cloud.

## Examples

### List all VPC networks
```sql
select network_id, name, description, created_at from yandexcloud_vpc_network;
```

### Find networks by name
```sql
select network_id, name from yandexcloud_vpc_network where name = 'default';
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| network_id  | text   | VPC network ID.                             |
| folder_id   | text   | Folder ID containing the network.           |
| name        | text   | Network name.                               |
| description | text   | Network description.                        |
| created_at  | text   | Network creation date (YYYY-MM-DD).         |
| labels      | jsonb  | Resource labels as key:value pairs.         | 