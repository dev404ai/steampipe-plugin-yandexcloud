select
  name,
  instance_id,
  resources ->> 'cores' as cores,
  resources ->> 'memory' as memory,
  boot_disk ->> 'diskId' as boot_disk_id,
  network_interfaces,
  labels,
  metadata,
  fqdn,
  hostname
from
  yandexcloud_compute_instance
where
  instance_id = '{{ output.resource_id.value }}'; 