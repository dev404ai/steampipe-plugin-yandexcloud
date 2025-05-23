---
title: Table: yandexcloud_compute_disk_placement_group
summary: Query information about Yandex Cloud Compute disk placement groups.
---

# Table: yandexcloud_compute_disk_placement_group

The `yandexcloud_compute_disk_placement_group` table allows you to query information about disk placement groups in Yandex Cloud Compute.

## Examples

### List all disk placement groups
```sql
select disk_placement_group_id, name, type, status, zone, created_at from yandexcloud_compute_disk_placement_group;
```

### Find disk placement groups by type
```sql
select disk_placement_group_id, name from yandexcloud_compute_disk_placement_group where type = 'spread';
```

## Columns
| Name                   | Type   | Description                                 |
|------------------------|--------|---------------------------------------------|
| disk_placement_group_id| text   | Disk placement group ID.                    |
| name                   | text   | Disk placement group name.                  |
| description            | text   | Disk placement group description.           |
| folder_id              | text   | Folder ID containing the placement group.   |
| zone                   | text   | Availability zone.                          |
| type                   | text   | Disk placement group type.                  |
| status                 | text   | Current status.                             |
| created_at             | text   | Creation date (YYYY-MM-DD).                 |
| labels                 | jsonb  | Resource labels as key:value pairs.         | 