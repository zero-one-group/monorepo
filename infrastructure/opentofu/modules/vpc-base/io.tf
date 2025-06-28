###################
###### INPUT #####
###################

variable "vpc_name" {
  type = string
}

variable "common_tags" {
  description = "Common tags"
}

variable "cidr_block" {
  type = string
}

variable "region_name" {
  type = string
}

###################
###### OUTPUT #####
###################

output "vpc_id" {
  value = aws_vpc.main.id
}
