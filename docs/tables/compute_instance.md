---
title: Table: yandexcloud_compute_instance
summary: Query information about Yandex Cloud Compute Instances.
---

# Table: yandexcloud_compute_instance

The `yandexcloud_compute_instance` table allows you to query information about virtual machine instances in Yandex Cloud.

## Examples

### List all compute instances
```sql
select id, name, status, zone, platform_id from yandexcloud_compute_instance;
```

### Find stopped instances
```sql
select id, name, status from yandexcloud_compute_instance where status = 'STOPPED';
```

## Columns
| Name                | Type   | Description                                 |
|---------------------|--------|---------------------------------------------|
| id                  | text   | Unique identifier for the instance.         |
| name                | text   | Name of the instance.                       |
| description         | text   | Description of the instance.                |
| folder_id           | text   | Folder ID where the instance resides.       |
| zone                | text   | Availability zone of the instance.          |
| platform_id         | text   | Platform type (e.g., standard-v1).          |
| status              | text   | Current status (e.g., RUNNING, STOPPED).    |
| created_at          | text   | Creation timestamp.                         |
| fqdn                | text   | Fully qualified domain name.                |
| hostname            | text   | Instance hostname.                          |
| service_account_id  | text   | Service account ID attached to the instance.|
| deletion_protection | bool   | Deletion protection flag.                   |
| labels              | jsonb  | Key-value labels assigned to the instance.  |
| resources           | jsonb  | Computing resources (CPU, RAM, GPU, etc.).  |
| metadata            | jsonb  | Instance metadata (e.g., ssh-keys).         |
| metadata_options    | jsonb  | Metadata access options.                    |
| boot_disk           | jsonb  | Boot disk information.                      |
| secondary_disks     | jsonb  | Secondary disks attached to the instance.   |
| network_interfaces  | jsonb  | Network interfaces attached to the instance.| 