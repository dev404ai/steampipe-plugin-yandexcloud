---
title: Table: yandexcloud_compute_placement_group
summary: Query information about Yandex Cloud Compute placement groups.
---

# Table: yandexcloud_compute_placement_group

The `yandexcloud_compute_placement_group` table allows you to query information about placement groups in Yandex Cloud Compute.

## Examples

### List all placement groups
```sql
select placement_group_id, name, type, status, zone, created_at from yandexcloud_compute_placement_group;
```

### Find placement groups by type
```sql
select placement_group_id, name from yandexcloud_compute_placement_group where type = 'spread';
```

## Columns
| Name               | Type   | Description                                 |
|--------------------|--------|---------------------------------------------|
| placement_group_id | text   | Placement group ID.                         |
| name               | text   | Placement group name.                       |
| description        | text   | Placement group description.                |
| folder_id          | text   | Folder ID containing the placement group.   |
| zone               | text   | Availability zone.                          |
| type               | text   | Placement group type.                       |
| status             | text   | Current status.                             |
| created_at         | text   | Placement group creation date (YYYY-MM-DD). |
| labels             | jsonb  | Resource labels as key:value pairs.         | 