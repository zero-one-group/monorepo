/**
 * # VPC
 *
 * This module references existing VPC and IGW
 *
 */

data "aws_vpc" "main" {
  filter {
    name   = "tag:Name"
    values = ["${var.vpc_name}-vpc"]
  }
}

data "aws_internet_gateway" "main" {
  filter {
    name   = "tag:Name"
    values = ["${var.vpc_name}-ig"]
  }

  filter {
    name   = "attachment.vpc-id"
    values = [data.aws_vpc.main.id]
  }
}
