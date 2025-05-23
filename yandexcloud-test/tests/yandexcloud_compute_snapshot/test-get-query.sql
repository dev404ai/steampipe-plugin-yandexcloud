-- To run this test, you need to set the snapshot_id variable in variables.json
select * from yandexcloud_compute_snapshot where snapshot_id = '{{ .snapshot_id }}'; 