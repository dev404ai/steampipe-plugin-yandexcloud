---
title: Yandex Cloud Plugin
summary: Query Yandex Cloud resources using SQL with Steampipe.
---

# Yandex Cloud Steampipe Plugin

The Yandex Cloud plugin for Steampipe enables you to query Yandex Cloud resources using SQL. Analyze your cloud infrastructure, audit security, and create dashboards with simple SQL queries.

## Features
- Query Compute Instances, Storage Buckets, Billing, IAM, and more
- Use SQL to join, filter, and aggregate cloud data
- Integrate with dashboards and reporting tools

## Getting Started
1. Install Steampipe: https://steampipe.io/downloads
2. Install this plugin: `steampipe plugin install dev404ai/yandexcloud`
3. Configure your credentials in `~/.steampipe/config/yandexcloud.spc` (see example in the config directory)

## Example Query
```sql
select name, status, zone from yandexcloud_compute_instance;
```

## Multiple Connections & Aggregator Examples

You can configure multiple Yandex Cloud connections in your Steampipe config to query across different accounts or environments. To aggregate data from several connections, use an aggregator connection. This allows you to run a single query across all specified connections as if they were one.

### Example: Multiple Connections
```hcl
connection "yandexcloud_01" {
  plugin      = "yandexcloud"
  service_account_key_file = "<YOUR_SERVICE_ACCOUNT_KEY_FILE_1>"
  cloud_id    = "<YOUR_CLOUD_ID_1>"
  folder_id   = "<YOUR_FOLDER_ID_1>"
}

connection "yandexcloud_02" {
  plugin      = "yandexcloud"
  service_account_key_file = "<YOUR_SERVICE_ACCOUNT_KEY_FILE_2>"
  cloud_id    = "<YOUR_CLOUD_ID_2>"
  folder_id   = "<YOUR_FOLDER_ID_2>"
}
```

### Example: Aggregator Connection
```hcl
connection "yandexcloud_all" {
  plugin      = "yandexcloud"
  type        = "aggregator"
  connections = ["yandexcloud_*"]
}
```

You can now run a query like this to get results from all configured Yandex Cloud connections:

```sql
select name, status, zone from yandexcloud_all.yandexcloud_compute_instance;
```

For more details, see the [Steampipe documentation on aggregators](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators).