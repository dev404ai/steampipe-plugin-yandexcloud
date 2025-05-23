select name, instance_id, status, zone, folder_id, description, platform_id, created_at
from yandexcloud_compute_instance
where instance_id = '{{ output.resource_id.value }}'; 