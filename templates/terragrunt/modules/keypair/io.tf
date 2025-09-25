variable "keyname" {
  type        = string
  description = "Keyname of keypair"
}

variable "public_key_path" {
  type        = string
  description = "Path file of public key for keypair"
}

variable "common_tags" {
  type        = map(any)
  description = "Common tags"
}

variable "region" {
  description = "AWS region"
  type        = string
}

output "keyname" {
  value       = aws_key_pair.key.key_name
  description = "Name of keypair"
}
