###################
###### INPUT #####
###################

variable "project_name" {
  type        = string
  description = "Project name for resource naming"
}

variable "environment" {
  type        = string
  description = "Environment name (dev, staging, prod)"
}

variable "region" {
  type        = string
  description = "AWS region"
}

variable "common_tags" {
  type        = map(string)
  description = "Common tags"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID the ALB is placed in"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Public subnet IDs the ALB is placed in"
}

variable "internal" {
  type        = bool
  description = "Whether the ALB is internal (false = internet-facing)"
  default     = false
}

variable "ingress_cidr_blocks" {
  type        = list(string)
  description = "CIDR blocks allowed to reach the ALB listeners"
  default     = ["0.0.0.0/0"]
}

variable "listener_port" {
  type        = number
  description = "HTTP listener port"
  default     = 80
}

variable "https_enabled" {
  type        = bool
  description = "Enable an additional HTTPS listener (needs a certificate, see modules/acm)"
  default     = false
}

variable "https_listener_port" {
  type        = number
  description = "HTTPS listener port"
  default     = 443
}

variable "certificate_arn" {
  type        = string
  description = "ARN of the ACM certificate for the HTTPS listener. Required when https_enabled = true."
  default     = null
}

variable "idle_timeout" {
  type        = number
  description = "ALB idle timeout in seconds"
  default     = 60
}

variable "services" {
  description = "Map of service key to target group + listener rule configuration. Keys must match the ECS service keys."
  type = map(object({
    container_port        = number
    alb_path_pattern      = optional(string, "/*")
    alb_health_check_path = optional(string, "/")
  }))
  default = {}
}

variable "access_logs_enabled" {
  type        = bool
  description = "Enable ALB access logs (requires an S3 bucket)"
  default     = false
}

variable "access_logs_bucket" {
  type        = string
  description = "S3 bucket name for ALB access logs. Required when access_logs_enabled = true."
  default     = null
}

###################
###### OUTPUT #####
###################

output "alb_arn" {
  description = "ARN of the ALB"
  value       = aws_lb.this.arn
}

output "alb_dns_name" {
  description = "DNS name of the ALB (route traffic here)"
  value       = aws_lb.this.dns_name
}

output "alb_security_group_id" {
  description = "Security group ID attached to the ALB"
  value       = aws_security_group.alb.id
}

output "tasks_security_group_id" {
  description = "Security group ID that allows traffic from the ALB into ECS tasks (attach to task services)"
  value       = aws_security_group.tasks.id
}

output "target_group_arns" {
  description = "Map of service key to target group ARN"
  value       = { for k, v in aws_lb_target_group.services : k => v.arn }
}

output "listener_arn" {
  description = "ARN of the default HTTP listener"
  value       = aws_lb_listener.http.arn
}