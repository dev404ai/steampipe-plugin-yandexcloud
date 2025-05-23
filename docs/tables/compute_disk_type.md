---
title: Table: yandexcloud_compute_disk_type
summary: Query information about Yandex Cloud Compute disk types.
---

# Table: yandexcloud_compute_disk_type

The `yandexcloud_compute_disk_type` table allows you to query information about disk types in Yandex Cloud Compute.

## Examples

### List all disk types
```sql
select disk_type_id, zone_id, name, description from yandexcloud_compute_disk_type;
```

### Find disk types by zone
```sql
select disk_type_id, name from yandexcloud_compute_disk_type where zone_id = 'ru-central1-a';
```

## Columns
| Name         | Type   | Description                                 |
|--------------|--------|---------------------------------------------|
| disk_type_id | text   | Disk type ID.                               |
| zone_id      | text   | Zone ID.                                    |
| name         | text   | Disk type name.                             |
| description  | text   | Disk type description.                      | 