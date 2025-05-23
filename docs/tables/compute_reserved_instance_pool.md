---
title: Table: yandexcloud_compute_reserved_instance_pool
summary: Query information about Yandex Cloud Compute reserved instance pools.
---

# Table: yandexcloud_compute_reserved_instance_pool

The `yandexcloud_compute_reserved_instance_pool` table allows you to query information about reserved instance pools in Yandex Cloud Compute.

## Examples

### List all reserved instance pools
```sql
select reserved_instance_pool_id, name, type, status, zone, created_at from yandexcloud_compute_reserved_instance_pool;
```

### Find reserved instance pools by type
```sql
select reserved_instance_pool_id, name from yandexcloud_compute_reserved_instance_pool where type = 'standard';
```

## Columns
| Name                     | Type   | Description                                 |
|--------------------------|--------|---------------------------------------------|
| reserved_instance_pool_id| text   | Reserved instance pool ID.                  |
| name                     | text   | Reserved instance pool name.                |
| description              | text   | Reserved instance pool description.         |
| folder_id                | text   | Folder ID containing the pool.              |
| zone                     | text   | Availability zone.                          |
| type                     | text   | Reserved instance pool type.                |
| status                   | text   | Current status.                             |
| created_at               | text   | Creation date (YYYY-MM-DD).                 |
| labels                   | jsonb  | Resource labels as key:value pairs.         | 