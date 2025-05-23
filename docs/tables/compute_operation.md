---
title: Table: yandexcloud_compute_operation
summary: Query information about Yandex Cloud Compute operations (get by ID only).
---

# Table: yandexcloud_compute_operation

The `yandexcloud_compute_operation` table allows you to query information about operations in Yandex Cloud Compute (get by operation ID only).

## Examples

### Get operation by ID
```sql
select operation_id, status, description, done, created_at from yandexcloud_compute_operation where operation_id = 'operation-123';
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| operation_id| text   | Operation ID.                               |
| folder_id   | text   | Folder ID.                                  |
| status      | text   | Operation status.                           |
| description | text   | Operation description.                      |
| done        | bool   | Operation done flag.                        |
| created_at  | text   | Operation creation date (YYYY-MM-DD).       |
| error       | jsonb  | Operation error details.                    |
| response    | jsonb  | Operation response details.                 |
| metadata    | jsonb  | Operation metadata.                         | 