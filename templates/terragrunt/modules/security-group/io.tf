###################
###### INPUT ######
###################

variable "security_groups" {
  type = map(object({
    name         = string
    description  = string
    ingress_rule = list(object({
      description                = optional(string, "Ingress Rule")
      from_port                  = optional(number, 0)
      to_port                    = optional(number, 0)
      protocol                   = optional(string, "-1")
      cidr_blocks                = optional(list(string), null)
      source_security_group_id   = optional(string, null)
      source_security_group_key  = optional(string, null)  # New field for self-references
    }))
    egress_rule = list(object({
      description                = optional(string, "Egress Rule")
      from_port                  = optional(number, 0)
      to_port                    = optional(number, 0)
      protocol                   = optional(string, "-1")
      cidr_blocks                = optional(list(string), ["0.0.0.0/0"])
      source_security_group_id   = optional(string, null)
      source_security_group_key  = optional(string, null)  # New field for self-references
    }))
  }))
  description = "Map of security groups to create with their rules"
}

variable "vpc_id" {
  type        = string
  description = "ID of VPC for security group"
}

variable "common_tags" {
  type        = map(string)
  description = "Common tags to apply to resources"
}

variable "region" {
  type        = string
  description = "AWS region"
}

###################
###### OUTPUT #####
###################

output "sg_ids" {
  value = {
    for key, sg in aws_security_group.sg : key => sg.id
  }
  description = "Map of Security Group IDs"
}

output "sg_arns" {
  value = {
    for key, sg in aws_security_group.sg : key => sg.arn
  }
  description = "Map of Security Group ARNs"
}

output "sg_names" {
  value = {
    for key, sg in aws_security_group.sg : key => sg.name
  }
  description = "Map of Security Group Names"
}

output "security_groups" {
  value = {
    for key, sg in aws_security_group.sg : key => {
      id   = sg.id
      arn  = sg.arn
      name = sg.name
    }
  }
  description = "Complete map of security groups with all details"
}
