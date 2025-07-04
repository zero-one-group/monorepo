###################
###### INPUT ######
###################

variable "sg_name" {
  type        = string
  description = "Name of security group"
}

variable "vpc_id" {
  type        = string
  description = "ID of VPC for security group"
}

variable "common_tags" {
  type        = map(any)
  description = "Common tags"
}

variable "ingress_rule" {
  type        = any
  default     = []
  description = "The ingress rule for security group"
}
variable "egress_rule" {
  type        = any
  default     = []
  description = "The egress rule for security group"
}

variable "description" {
  type        = string
  description = "Description for security group"
}
###################
###### OUTPUT #####
###################
output "sg_id" {
  value       = aws_security_group.sg.id
  description = "ID of Security Group"
}
