---
title: Table: yandexcloud_compute_snapshot
summary: Query information about Yandex Cloud Compute disk snapshots.
---

# Table: yandexcloud_compute_snapshot

The `yandexcloud_compute_snapshot` table allows you to query information about disk snapshots in Yandex Cloud Compute.

## Examples

### List all snapshots
```sql
select snapshot_id, name, status, created_at from yandexcloud_compute_snapshot;
```

### Find snapshots by status
```sql
select snapshot_id, name from yandexcloud_compute_snapshot where status = 'READY';
```

## Columns
| Name           | Type   | Description                                 |
|----------------|--------|---------------------------------------------|
| snapshot_id    | text   | Snapshot ID.                               |
| name           | text   | Snapshot name.                             |
| description    | text   | Snapshot description.                      |
| folder_id      | text   | Folder ID containing the snapshot.         |
| zone           | text   | Availability zone.                         |
| status         | text   | Current status.                            |
| created_at     | text   | Snapshot creation date (YYYY-MM-DD).       |
| source_disk_id | text   | Source disk ID.                            |
| size           | text   | Snapshot size (bytes).                     | 