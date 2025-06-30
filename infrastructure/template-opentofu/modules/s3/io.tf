# Variables
variable "bucket_name" {
  description = "Name of the S3 bucket"
  type        = string
}

variable "enable_versioning" {
  description = "Enable versioning on the bucket"
  type        = bool
  default     = false
}

variable "enable_lifecycle_rule" {
  description = "Enable lifecycle rules on the bucket"
  type        = bool
  default     = false
}

variable "transition_days" {
  description = "Number of days before transitioning objects to STANDARD_IA"
  type        = number
  default     = 30
}

variable "expiration_days" {
  description = "Number of days before expiring objects"
  type        = number
  default     = 90
}

variable "tags" {
  description = "Tags to be applied to the bucket"
  type        = map(string)
  default     = {}
}

variable "enable_lb_logging" {
  description = "Enable the bucket for load balancer logging"
  type        = bool
  default     = false
}

variable "access_logs_prefix" {
  description = "Prefix for load balancer logs in the bucket"
  type        = string
  default     = ""
}

# Outputs
output "bucket_id" {
  description = "The name of the bucket"
  value       = aws_s3_bucket.this.id
}

output "bucket_arn" {
  description = "The ARN of the bucket"
  value       = aws_s3_bucket.this.arn
}

output "bucket_domain_name" {
  description = "The bucket domain name"
  value       = aws_s3_bucket.this.bucket_domain_name
}
