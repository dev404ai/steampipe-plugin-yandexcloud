# Example Steampipe connection configuration for Yandex Cloud
# See: https://hub.steampipe.io/plugins/dev404ai/yandexcloud

connection "yandexcloud" {
  plugin = "yandexcloud"

  # Path to the service account key file (JSON)
  # service_account_key_file = "/path/to/key.json"

  # Yandex Cloud cloud ID (required)
  cloud_id = "b1g7xxxxxx"

  # Yandex Cloud folder ID (required)
  folder_id = "b1g7yyyyyy"

  # Log level: error, info, or debug (optional)
  # log_level = "info"
} 