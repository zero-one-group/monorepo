# Variables
variable "buckets" {
  description = "Map of bucket configurations"
  type = map(object({
    bucket_name            = string
    enable_versioning      = optional(bool, false)
    enable_lifecycle_rule  = optional(bool, false)
    transition_days        = optional(number, 30)
    expiration_days        = optional(number, 90)
    enable_lb_logging      = optional(bool, false)
    access_logs_prefix     = optional(string, "")
    cors_origin           = optional(list(string), [])
    enable_public_read     = optional(bool, false)
    custom_policy         = optional(string, "")
    allow_public_policy   = optional(bool, false)
  }))
}

variable "common_tags" {
  description = "Common tags to be applied to the bucket"
  type        = map(any)
}

variable "region" {
  description = "AWS region"
  type        = string
}

# Outputs
output "bucket_ids" {
  description = "Map of bucket names"
  value       = { for k, v in aws_s3_bucket.this : k => v.id }
}

output "bucket_arns" {
  description = "Map of bucket ARNs"
  value       = { for k, v in aws_s3_bucket.this : k => v.arn }
}

output "bucket_domain_names" {
  description = "Map of bucket domain names"
  value       = { for k, v in aws_s3_bucket.this : k => v.bucket_domain_name }
}
