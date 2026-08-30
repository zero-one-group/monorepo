variable "region" {
  description = "AWS region (required by the generated provider config, even while this module is a placeholder)"
  type        = string
}

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
#   cluster_identifier
#   engine / engine_version
#   master_username / master_password
#   database_name
#   instance_count
#   instance_class
#   storage / storage_encrypted
#   backup_retention_period
#   maintenance_window
#   enable_deletion_protection
#
# Planned outputs:
#   cluster_arn
#   endpoint
#   port
#   security_group_id
#   parameter_group_id

terraform {
  required_version = ">= 1.0"
}
