select
  security_group_id,
  name,
  network_id,
  folder_id
from
  yandexcloud_vpc_security_group
limit 2; 