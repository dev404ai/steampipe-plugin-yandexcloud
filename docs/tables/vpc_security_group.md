---
title: Table: yandexcloud_vpc_security_group
summary: Query information about Yandex Cloud VPC security groups.
---

# Table: yandexcloud_vpc_security_group

The `yandexcloud_vpc_security_group` table allows you to query information about security groups in Yandex Cloud VPC.

## Examples

### List all security groups
```sql
select security_group_id, name, network_id, created_at from yandexcloud_vpc_security_group;
```

### Find security groups by network
```sql
select security_group_id, name from yandexcloud_vpc_security_group where network_id = 'network-123';
```

## Columns
| Name              | Type   | Description                                 |
|-------------------|--------|---------------------------------------------|
| security_group_id | text   | Security group ID.                          |
| folder_id         | text   | Folder ID containing the security group.     |
| network_id        | text   | Network ID to which the security group belongs. |
| name              | text   | Security group name.                        |
| description       | text   | Security group description.                 |
| created_at        | text   | Security group creation date (YYYY-MM-DD).  |
| labels            | jsonb  | Resource labels as key:value pairs.         |
| rules             | jsonb  | Security group rules.                       | 