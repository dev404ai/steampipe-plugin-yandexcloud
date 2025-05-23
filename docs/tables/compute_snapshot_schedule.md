---
title: Table: yandexcloud_compute_snapshot_schedule
summary: Query information about Yandex Cloud Compute snapshot schedules.
---

# Table: yandexcloud_compute_snapshot_schedule

The `yandexcloud_compute_snapshot_schedule` table allows you to query information about snapshot schedules in Yandex Cloud Compute.

## Examples

### List all snapshot schedules
```sql
select snapshot_schedule_id, name, status, zone, created_at from yandexcloud_compute_snapshot_schedule;
```

### Find snapshot schedules by status
```sql
select snapshot_schedule_id, name from yandexcloud_compute_snapshot_schedule where status = 'ACTIVE';
```

## Columns
| Name                 | Type   | Description                                 |
|----------------------|--------|---------------------------------------------|
| snapshot_schedule_id | text   | Snapshot schedule ID.                       |
| name                 | text   | Snapshot schedule name.                     |
| description          | text   | Snapshot schedule description.              |
| folder_id            | text   | Folder ID containing the schedule.          |
| zone                 | text   | Availability zone.                          |
| status               | text   | Current status.                             |
| created_at           | text   | Creation date (YYYY-MM-DD).                 |
| labels               | jsonb  | Resource labels as key:value pairs.         | 