---
title: Table: yandexcloud_vpc_route_table
summary: Query information about Yandex Cloud VPC route tables.
---

# Table: yandexcloud_vpc_route_table

The `yandexcloud_vpc_route_table` table allows you to query information about route tables in Yandex Cloud VPC.

## Examples

### List all route tables
```sql
select route_table_id, name, network_id, created_at from yandexcloud_vpc_route_table;
```

### Find route tables by network
```sql
select route_table_id, name from yandexcloud_vpc_route_table where network_id = 'network-123';
```

## Columns
| Name           | Type   | Description                                 |
|----------------|--------|---------------------------------------------|
| route_table_id | text   | VPC route table ID.                         |
| folder_id      | text   | Folder ID containing the route table.        |
| network_id     | text   | Network ID to which the route table belongs. |
| name           | text   | Route table name.                            |
| description    | text   | Route table description.                     |
| created_at     | text   | Route table creation date (YYYY-MM-DD).      |
| labels         | jsonb  | Resource labels as key:value pairs.          |
| static_routes  | jsonb  | List of static routes.                       | 