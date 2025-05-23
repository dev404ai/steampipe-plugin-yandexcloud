---
title: Table: yandexcloud_compute_host_group
summary: Query information about Yandex Cloud Compute host groups.
---

# Table: yandexcloud_compute_host_group

The `yandexcloud_compute_host_group` table allows you to query information about host groups in Yandex Cloud Compute.

## Examples

### List all host groups
```sql
select host_group_id, name, type, status, zone, created_at from yandexcloud_compute_host_group;
```

### Find host groups by type
```sql
select host_group_id, name from yandexcloud_compute_host_group where type = 'dedicated';
```

## Columns
| Name           | Type   | Description                                 |
|----------------|--------|---------------------------------------------|
| host_group_id  | text   | Host group ID.                              |
| name           | text   | Host group name.                            |
| description    | text   | Host group description.                     |
| folder_id      | text   | Folder ID containing the host group.        |
| zone           | text   | Availability zone.                          |
| type           | text   | Host group type.                            |
| status         | text   | Current status.                             |
| created_at     | text   | Host group creation date (YYYY-MM-DD).      |
| labels         | jsonb  | Resource labels as key:value pairs.         | 