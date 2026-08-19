variable "enabled" {
  description = "Set to true once the module is implemented"
  type        = bool
  default     = false
}

# Planned inputs (module not implemented yet - structure only):
#   project_name
#   environment
#   region
#   common_tags
#   vpc_id
#   subnet_ids
#   security_group_ids
#   db_identifier
#   engine / engine_version
#   username / password
#   database_name
#   instance_class
#   allocated_storage / max_allocated_storage
#   storage_encrypted
#   backup_retention_period
#   maintenance_window
#   deletion_protection
#
# Planned outputs:
#   db_instance_arn
#   endpoint
#   port
#   security_group_id

terraform {
  required_version = ">= 1.0"
}
