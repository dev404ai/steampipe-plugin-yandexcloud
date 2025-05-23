---
title: Table: yandexcloud_compute_host_type
summary: Query information about Yandex Cloud Compute host types.
---

# Table: yandexcloud_compute_host_type

The `yandexcloud_compute_host_type` table allows you to query information about host types in Yandex Cloud Compute.

## Examples

### List all host types
```sql
select host_type_id, zone_id, name, description from yandexcloud_compute_host_type;
```

### Find host types by zone
```sql
select host_type_id, name from yandexcloud_compute_host_type where zone_id = 'ru-central1-a';
```

## Columns
| Name         | Type   | Description                                 |
|--------------|--------|---------------------------------------------|
| host_type_id | text   | Host type ID.                               |
| zone_id      | text   | Zone ID.                                    |
| name         | text   | Host type name.                             |
| description  | text   | Host type description.                      | 