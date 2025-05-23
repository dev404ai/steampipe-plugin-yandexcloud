---
title: Table: yandexcloud_compute_gpu_cluster
summary: Query information about Yandex Cloud Compute GPU clusters.
---

# Table: yandexcloud_compute_gpu_cluster

The `yandexcloud_compute_gpu_cluster` table allows you to query information about GPU clusters in Yandex Cloud Compute.

## Examples

### List all GPU clusters
```sql
select gpu_cluster_id, name, type, status, zone, created_at from yandexcloud_compute_gpu_cluster;
```

### Find GPU clusters by type
```sql
select gpu_cluster_id, name from yandexcloud_compute_gpu_cluster where type = 'standard';
```

## Columns
| Name           | Type   | Description                                 |
|----------------|--------|---------------------------------------------|
| gpu_cluster_id | text   | GPU cluster ID.                             |
| name           | text   | GPU cluster name.                           |
| description    | text   | GPU cluster description.                    |
| folder_id      | text   | Folder ID containing the GPU cluster.       |
| zone           | text   | Availability zone.                          |
| type           | text   | GPU cluster type.                           |
| status         | text   | Current status.                             |
| created_at     | text   | GPU cluster creation date (YYYY-MM-DD).     |
| labels         | jsonb  | Resource labels as key:value pairs.         | 