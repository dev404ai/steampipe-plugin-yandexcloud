---
title: Table: yandexcloud_vpc_operation
summary: Query information about Yandex Cloud VPC operations.
---

# Table: yandexcloud_vpc_operation

The `yandexcloud_vpc_operation` table allows you to query information about operations in Yandex Cloud VPC.

## Examples

### List all VPC operations
```sql
select operation_id, description, created_by, done, created_at from yandexcloud_vpc_operation;
```

### Find completed operations
```sql
select operation_id, description from yandexcloud_vpc_operation where done = true;
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| operation_id| text   | Operation ID.                               |
| description | text   | Operation description.                      |
| created_at  | text   | Operation creation date (YYYY-MM-DD).       |
| created_by  | text   | ID of the user or service account who initiated the operation. |
| modified_at | text   | The time when the operation was last modified. |
| done        | bool   | If true, the operation is completed.         |
| metadata    | jsonb  | Service-specific metadata associated with the operation. | 