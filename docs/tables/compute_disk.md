---
title: Table: yandexcloud_compute_disk
summary: Query information about Yandex Cloud Compute disks.
---

# Table: yandexcloud_compute_disk

The `yandexcloud_compute_disk` table allows you to query information about disks in Yandex Cloud Compute.

## Examples

### List all disks
```sql
select disk_id, name, type_id, size, status, zone_id from yandexcloud_compute_disk;
```

### Find disks by type
```sql
select disk_id, name from yandexcloud_compute_disk where type_id = 'network-ssd';
```

## Columns
| Name         | Type   | Description                                 |
|--------------|--------|---------------------------------------------|
| disk_id      | text   | Disk ID.                                    |
| name         | text   | Disk name.                                  |
| description  | text   | Disk description.                           |
| folder_id    | text   | Folder ID containing the disk.              |
| type_id      | text   | Disk type ID.                               |
| size         | text   | Disk size (bytes).                          |
| status       | text   | Current status.                             |
| zone_id      | text   | Availability zone.                          |
| created_at   | text   | Disk creation date (YYYY-MM-DD).            |
| source_image_id | text| Source image ID.                            |
| labels       | jsonb  | Resource labels as key:value pairs.         | 