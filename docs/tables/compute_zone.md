---
title: Table: yandexcloud_compute_zone
summary: Query information about Yandex Cloud Compute zones.
---

# Table: yandexcloud_compute_zone

The `yandexcloud_compute_zone` table allows you to query information about availability zones in Yandex Cloud Compute.

## Examples

### List all compute zones
```sql
select zone_id, region_id, name, status from yandexcloud_compute_zone;
```

### Find zones by status
```sql
select zone_id, name from yandexcloud_compute_zone where status = 'UP';
```

## Columns
| Name      | Type   | Description                                 |
|-----------|--------|---------------------------------------------|
| zone_id   | text   | Zone ID.                                    |
| region_id | text   | Region ID.                                  |
| name      | text   | Zone name.                                  |
| status    | text   | Zone status.                                | 