variable "instance_name" {
  type        = string
  description = "Name of instance"
}

variable "ami_instance" {
  type        = string
  description = "AMI id for instance"
}

variable "keyname" {
  type        = string
  description = "Key name for instance, get from keypair module"
}

variable "security_group" {
  type        = list(any)
  description = "Security group for instance"
}

variable "subnet" {
  type        = string
  description = "Subnet for instance"
}

variable "common_tags" {
  type        = map(any)
  description = "Common tags"
}

variable "instance_type" {
  description = "Instance type of ec2 instance"
  type        = string
}

variable "user_data" {
  description = "Path file and vars of template user-data ec2"
  type = object({
    filename = string
    vars     = optional(map(string), { aws_region = "" })
  })
}

variable "root_block_device" {
  description = "Configuration for the root block device"
  type = object({
    volume_type           = string
    volume_size           = number
    delete_on_termination = bool
    encrypted             = bool
    kms_key_id            = string
  })
}

variable "hostname" {
  description = "Custom hostname for the instance"
  type        = string
  default     = "" # Default empty in case not all instances need custom hostname
}

output "ec2_public_ip" {
  value       = aws_instance.instance.public_ip
  description = "IP public of ec2 instance"
}

output "ec2_private_ip" {
  value       = aws_instance.instance.private_ip
  description = "IP private of ec2 instance"
}

output "ec2_instance_id" {
  value       = aws_instance.instance.id
  description = "Instance ID of ec2 instance"
}

output "instance_public_dns" {
  value = aws_instance.instance.public_dns
}
