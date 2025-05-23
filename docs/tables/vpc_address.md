---
title: Table: yandexcloud_vpc_address
summary: Query information about Yandex Cloud VPC addresses.
---

# Table: yandexcloud_vpc_address

The `yandexcloud_vpc_address` table allows you to query information about VPC addresses in Yandex Cloud.

## Examples

### List all VPC addresses
```sql
select address_id, name, type, ip_version, reserved, used, created_at from yandexcloud_vpc_address;
```

### Find reserved addresses
```sql
select address_id, name from yandexcloud_vpc_address where reserved = true;
```

## Columns
| Name               | Type   | Description                                 |
|--------------------|--------|---------------------------------------------|
| address_id         | text   | VPC address ID.                             |
| folder_id          | text   | Folder ID containing the address.           |
| created_at         | text   | Address creation date (YYYY-MM-DD).         |
| name               | text   | Address name.                               |
| description        | text   | Address description.                        |
| labels             | jsonb  | Resource labels as key:value pairs.         |
| external_ipv4_address | jsonb| External IPv4 address specification.        |
| reserved           | bool   | Specifies if address is reserved.           |
| used               | bool   | Specifies if address is used.               |
| type               | text   | Type of the IP address (INTERNAL/EXTERNAL). |
| ip_version         | text   | Version of the IP address (IPV4/IPV6).      |
| deletion_protection| bool   | Specifies if address is protected from deletion. |
| dns_records        | jsonb  | DNS record specifications.                  | 