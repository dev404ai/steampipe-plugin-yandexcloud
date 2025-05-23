---
title: Table: yandexcloud_billing_resource_usage
summary: Query billing resource usage in Yandex Cloud.
---

# Table: yandexcloud_billing_resource_usage

The `yandexcloud_billing_resource_usage` table allows you to query billing and usage data for Yandex Cloud resources.

## Examples

### List all billing resource usage records
```sql
select resource_id, product, usage, cost, usage_date from yandexcloud_billing_resource_usage;
```

### Find usage for a specific product
```sql
select resource_id, usage, cost from yandexcloud_billing_resource_usage where product = 'Compute Cloud';
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| resource_id | text   | Unique identifier for the resource.         |
| product     | text   | Name of the product or service.             |
| usage       | double | Amount of resource used.                    |
| cost        | double | Cost incurred for the usage.                |
| usage_date  | text   | Date of the usage record.                   |
| folder_id   | text   | Folder ID associated with the usage.        | 