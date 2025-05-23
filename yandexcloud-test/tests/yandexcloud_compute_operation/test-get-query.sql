-- Replace 'operation-123' with the actual operation_id for the test
select
  operation_id,
  status,
  folder_id
from
  yandexcloud_compute_operation
where
  operation_id = 'operation-123'; 