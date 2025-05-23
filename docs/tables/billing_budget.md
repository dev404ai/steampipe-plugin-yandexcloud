---
title: Table: yandexcloud_billing_budget
summary: Query information about Yandex Cloud Billing Budgets.
---

# Table: yandexcloud_billing_budget

The `yandexcloud_billing_budget` table allows you to query information about billing budgets in Yandex Cloud.

## Examples

### List all billing budgets
```sql
select budget_id, name, amount, created_at from yandexcloud_billing_budget;
```

## Columns
| Name        | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| budget_id   | text   | Budget ID.                                  |
| name        | text   | Budget name.                                |
| amount      | text   | Budget amount.                              |
| created_at  | text   | Creation date (RFC3339).                    | 