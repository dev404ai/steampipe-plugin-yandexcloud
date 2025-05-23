---
title: Table: yandexcloud_billing_account
summary: Query information about Yandex Cloud Billing Accounts.
---

# Table: yandexcloud_billing_account

The `yandexcloud_billing_account` table allows you to query information about billing accounts in Yandex Cloud.

## Examples

### List all billing accounts
```sql
select id, name, active, created_at from yandexcloud_billing_account;
```

### Find active billing accounts
```sql
select id, name from yandexcloud_billing_account where active = true;
```

## Columns
| Name         | Type   | Description                                 |
|--------------|--------|---------------------------------------------|
| id           | text   | Billing account ID.                         |
| name         | text   | Account name.                               |
| created_at   | text   | Creation date (RFC3339).                    |
| country_code | text   | Country code.                               |
| balance      | text   | Current balance.                            |
| currency     | text   | Currency code.                              |
| active       | bool   | Is account active?                          |
| labels       | jsonb  | Resource labels as key:value pairs.         | 