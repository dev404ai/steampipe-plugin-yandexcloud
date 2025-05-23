-- To run this test, you need to set the image_id variable in variables.json
select * from yandexcloud_compute_image where image_id = '{{ .image_id }}'; 