 variable "folder_id" {
  description = "Yandex Cloud folder ID for test resources."
  type        = string
}

variable "zone" {
  description = "Yandex Cloud zone for test instance."
  type        = string
  default     = "ru-central1-a"
}

variable "instance_name" {
  description = "Name for the test compute instance."
  type        = string
  default     = "steampipe-test-instance"
}

variable "platform_id" {
  description = "Platform ID for the test instance."
  type        = string
  default     = "standard-v3"
}

variable "cores" {
  description = "Number of vCPUs for the test instance."
  type        = number
  default     = 2
}

variable "memory" {
  description = "RAM size in GB for the test instance."
  type        = number
  default     = 4
}

variable "labels" {
  description = "Labels for the test instance."
  type        = map(string)
  default     = {
    "env" = "test"
    "project" = "steampipe"
  }
}
