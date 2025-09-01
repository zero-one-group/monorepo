variable "common_tags" {
  type        = map(any)
  description = "Common tags"
}

variable "repo_name" {
  type        = string
  description = "Name of repository"
}

variable "untagged_image_expiration_days" {
  type        = number
  description = "Number of days after which untagged images will be removed"
  default     = 1
}

variable "region" {
  description = "AWS region"
  type        = string
}


output "repository_url" {
  value       = aws_ecr_repository.repo.repository_url
  description = "URL of Repository"
}
