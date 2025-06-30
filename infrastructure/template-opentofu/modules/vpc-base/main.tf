/**
 * # VPC
 *
 * This module only make vpc and igw so it can reuse on vpc module
 *
 */

resource "aws_vpc" "main" {
  cidr_block           = var.cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-vpc"
    })
  )
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-ig"
    })
  )
}
