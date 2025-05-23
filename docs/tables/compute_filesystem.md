---
title: Table: yandexcloud_compute_filesystem
summary: Query information about Yandex Cloud Compute filesystems.
---

# Table: yandexcloud_compute_filesystem

The `yandexcloud_compute_filesystem` table allows you to query information about filesystems in Yandex Cloud Compute.

## Examples

### List all filesystems
```sql
select filesystem_id, name, type_id, size, status, zone_id from yandexcloud_compute_filesystem;
```

### Find filesystems by type
```sql
select filesystem_id, name from yandexcloud_compute_filesystem where type_id = 'network-ssd';
```

## Columns
| Name         | Type   | Description                                 |
|--------------|--------|---------------------------------------------|
| filesystem_id| text   | Filesystem ID.                              |
| name         | text   | Filesystem name.                            |
| description  | text   | Filesystem description.                     |
| folder_id    | text   | Folder ID containing the filesystem.        |
| type_id      | text   | Filesystem type ID.                         |
| size         | text   | Filesystem size (bytes).                    |
| status       | text   | Current status.                             |
| zone_id      | text   | Availability zone.                          |
| created_at   | text   | Filesystem creation date (YYYY-MM-DD).      |
| labels       | jsonb  | Resource labels as key:value pairs.         | 