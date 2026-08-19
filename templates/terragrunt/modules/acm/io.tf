variable "enabled" {
  description = "Set to true once the module is implemented"
  type        = bool
  default     = false
}

# Planned inputs (module not implemented yet - structure only):
#   project_name
#   environment
#   region
#   domain_name
#   alternative_names
#   validate_options
#   tags
#
# Planned outputs:
#   certificate_arn

terraform {
  required_version = ">= 1.0"
}
