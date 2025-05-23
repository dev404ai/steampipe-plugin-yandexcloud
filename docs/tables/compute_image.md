---
title: Table: yandexcloud_compute_image
summary: Query information about Yandex Cloud Compute disk images.
---

# Table: yandexcloud_compute_image

The `yandexcloud_compute_image` table allows you to query information about disk images in Yandex Cloud Compute.

## Examples

### List all images
```sql
select image_id, name, status, created_at from yandexcloud_compute_image;
```

### Find images by family
```sql
select image_id, name from yandexcloud_compute_image where family = 'ubuntu-2004-lts';
```

## Columns
| Name         | Type   | Description                                 |
|--------------|--------|---------------------------------------------|
| image_id     | text   | Image ID.                                   |
| name         | text   | Image name.                                 |
| description  | text   | Image description.                          |
| folder_id    | text   | Folder ID containing the image.             |
| family       | text   | Image family.                               |
| product_ids  | jsonb  | Product IDs associated with the image.      |
| status       | text   | Current status.                             |
| created_at   | text   | Image creation date (YYYY-MM-DD).           |
| min_disk_size| text   | Minimum disk size required (bytes).         | 