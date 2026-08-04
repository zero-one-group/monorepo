###################
###### INPUT ######
###################

variable "project_name" {
  type        = string
  description = "Project name for resource naming"
}

variable "region" {
  type        = string
  description = "AWS region"
}

variable "aws_account_id" {
  type        = string
  description = "AWS account ID"
}

variable "common_tags" {
  type        = map(any)
  description = "Common tags"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID tasks run in"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Subnet IDs tasks are placed in"
}

variable "security_group_ids" {
  type        = list(string)
  description = "Security group IDs attached to tasks"
  default     = []
}

variable "ecr_repository_url" {
  type        = string
  description = "ECR repository URL images are pulled from"
}

variable "image_tag" {
  type        = string
  description = "Image tag prefixed to the image reference for each entry: <ecr_repository_url>:<image_tag>-<key> (matches the <tag>-<app> tag pushed by the build pipeline). Override per entry with the `image` field."
}

variable "dns_namespace_name" {
  type        = string
  description = "Private DNS namespace (Cloud Map) used for service discovery between ECS tasks, e.g. 'myapp.local'. Only created if at least one service has enable_service_discovery = true. Defaults to '<project_name>.local'."
  default     = null
}

variable "services" {
  description = "Always-on workloads: task definition + ECS service. Optional service discovery and/or ALB per service (set enable_alb = true with container_port)."
  type = map(object({
    image                    = optional(string)
    command                  = optional(list(string), [])
    cpu                      = number
    memory                   = number
    environment              = optional(map(string), {})
    assign_public_ip         = optional(bool, false)
    enable_service_discovery = optional(bool, false)
    enable_alb               = optional(bool, false)
    container_port           = optional(number)
    alb_path_pattern         = optional(string, "/*")
    alb_health_check_path    = optional(string, "/")
    desired_count            = optional(number, 1)
    s3_bucket_arns           = optional(list(string), [])
    ssm_parameter_paths      = optional(list(string), [])
  }))
  default = {}
}

variable "jobs" {
  description = "On-demand workloads: task definition only, launched via ecs:RunTask. No ECS service is created."
  type = map(object({
    image               = optional(string)
    command             = optional(list(string), [])
    cpu                 = number
    memory              = number
    environment         = optional(map(string), {})
    assign_public_ip    = optional(bool, false)
    s3_bucket_arns      = optional(list(string), [])
    ssm_parameter_paths = optional(list(string), [])
  }))
  default = {}
}

variable "alb_internal" {
  type        = bool
  description = "Whether the ALB is internal (false = internet-facing)"
  default     = false
}

variable "alb_subnet_ids" {
  type        = list(string)
  description = "Subnet IDs the ALB is placed in. Required when any service has enable_alb = true."
  default     = []
}

variable "alb_ingress_cidr_blocks" {
  type        = list(string)
  description = "CIDR blocks allowed to reach the ALB listener port"
  default     = ["0.0.0.0/0"]
}

variable "alb_listener_port" {
  type        = number
  description = "ALB HTTP listener port"
  default     = 80
}

variable "alb_idle_timeout" {
  type        = number
  description = "ALB idle timeout in seconds"
  default     = 60
}

###################
###### OUTPUT #####
###################

output "cluster_arn" {
  description = "ARN of the ECS cluster"
  value       = aws_ecs_cluster.this.arn
}

output "cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.this.name
}

output "service_task_definition_arns" {
  description = "Map of service key to task definition ARN"
  value       = { for k, v in aws_ecs_task_definition.services : k => v.arn }
}

output "job_task_definition_arns" {
  description = "Map of job key to task definition ARN (family:revision), for use with ecs:RunTask"
  value       = { for k, v in aws_ecs_task_definition.jobs : k => v.arn }
}

output "job_task_definition_families" {
  description = "Map of job key to task definition family name"
  value       = { for k, v in aws_ecs_task_definition.jobs : k => v.family }
}

output "service_discovery_dns_names" {
  description = "Map of service key to its Cloud Map DNS name (only for services with enable_service_discovery = true)"
  value       = { for k in keys(aws_service_discovery_service.this) : k => "${k}.${local.dns_namespace_name}" }
}

output "alb_arn" {
  description = "ARN of the ALB. Empty unless at least one service has enable_alb = true."
  value       = local.alb_enabled ? aws_lb.this[0].arn : null
}

output "alb_dns_name" {
  description = "DNS name of the ALB (route traffic here). Empty unless at least one service has enable_alb = true."
  value       = local.alb_enabled ? aws_lb.this[0].dns_name : null
}

output "alb_security_group_id" {
  description = "ID of the ALB security group. Empty unless at least one service has enable_alb = true."
  value       = local.alb_enabled ? aws_security_group.alb[0].id : null
}

output "alb_target_group_arns" {
  description = "Map of service key to target group ARN for services with enable_alb = true"
  value       = { for k, v in aws_lb_target_group.services : k => v.arn }
}
