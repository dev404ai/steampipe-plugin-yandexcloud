---
title: Table: yandexcloud_vpc_gateway
summary: Query information about Yandex Cloud VPC gateways.
---

# Table: yandexcloud_vpc_gateway

The `yandexcloud_vpc_gateway` table allows you to query information about VPC gateways in Yandex Cloud.

## Examples

### List all VPC gateways
```sql
select gateway_id, name, created_at from yandexcloud_vpc_gateway;
```

### Find gateways by name
```sql
select gateway_id, name from yandexcloud_vpc_gateway where name = 'default';
```

## Columns
| Name                  | Type   | Description                                 |
|-----------------------|--------|---------------------------------------------|
| gateway_id            | text   | VPC gateway ID.                             |
| folder_id             | text   | Folder ID containing the gateway.           |
| created_at            | text   | Gateway creation date (YYYY-MM-DD).         |
| name                  | text   | Gateway name.                               |
| description           | text   | Gateway description.                        |
| labels                | jsonb  | Resource labels as key:value pairs.         |
| shared_egress_gateway | jsonb  | Shared egress gateway specification.        | 