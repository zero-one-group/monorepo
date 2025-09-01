/**
 * # VPC
 *
 * This module provisions subnets for a specific environment
 * using existing VPC and Internet Gateway from vpc-base.
 */

# Reference the shared VPC
data "aws_vpc" "main" {
  id = var.vpc_id
}

# Reference the shared Internet Gateway
data "aws_internet_gateway" "main" {
  internet_gateway_id = var.internet_gateway_id
}

# Get available AZs
data "aws_availability_zones" "available" {
  state = "available"
}

#####################################################
# Public Subnets                                   #
#####################################################
resource "aws_subnet" "public_a" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, var.subnet_offset + 1)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[0]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public-a"
    })
  )
}

resource "aws_subnet" "public_b" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, var.subnet_offset + 2)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[1]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public-b"
    })
  )
}

resource "aws_subnet" "public_c" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, var.subnet_offset + 3)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[2]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public-c"
    })
  )
}

# Public Route Table
resource "aws_route_table" "public" {
  vpc_id = data.aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public"
    })
  )
}

resource "aws_route" "public_internet_access" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = data.aws_internet_gateway.main.id
}

# Public Route Table Associations
resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_b" {
  subnet_id      = aws_subnet.public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_c" {
  subnet_id      = aws_subnet.public_c.id
  route_table_id = aws_route_table.public.id
}
