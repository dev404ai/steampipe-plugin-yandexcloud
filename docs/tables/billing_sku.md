---
title: Table: yandexcloud_billing_sku
summary: Query information about Yandex Cloud Billing SKUs.
---

# Table: yandexcloud_billing_sku

The `yandexcloud_billing_sku` table allows you to query information about billing SKUs in Yandex Cloud.

## Examples

### List all billing SKUs
```sql
select sku_id, name, description from yandexcloud_billing_sku;
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| sku_id      | text   | SKU ID.                                     |
| name        | text   | SKU name.                                   |
| description | text   | SKU description.                            | 