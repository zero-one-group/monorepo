/**
 * # VPC
 *
 * This module creates a complete VPC infrastructure for a specific environment
 * including VPC, Internet Gateway, subnets, and routing.
 * Security groups are managed in a separate module.
 */

# Get available AZs
data "aws_availability_zones" "available" {
  state = "available"
}

#####################################################
# VPC                                               #
#####################################################
resource "aws_vpc" "main" {
  cidr_block           = var.cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-vpc"
      "Environment" = var.environment
    })
  )
}

#####################################################
# Internet Gateway                                  #
#####################################################
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-igw"
      "Environment" = var.environment
    })
  )
}

#####################################################
# Public Subnets                                   #
#####################################################
resource "aws_subnet" "public_a" {
  cidr_block              = cidrsubnet(var.cidr_block, 8, 1)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[0]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-public-a"
      "Environment" = var.environment
      "Type"        = "public"
    })
  )
}

resource "aws_subnet" "public_b" {
  cidr_block              = cidrsubnet(var.cidr_block, 8, 2)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[1]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-public-b"
      "Environment" = var.environment
      "Type"        = "public"
    })
  )
}

resource "aws_subnet" "public_c" {
  cidr_block              = cidrsubnet(var.cidr_block, 8, 3)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[2]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-public-c"
      "Environment" = var.environment
      "Type"        = "public"
    })
  )
}

#####################################################
# Private Subnets                                  #
#####################################################
resource "aws_subnet" "private_a" {
  cidr_block        = cidrsubnet(var.cidr_block, 8, 4)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-private-a"
      "Environment" = var.environment
      "Type"        = "private"
    })
  )
}

resource "aws_subnet" "private_b" {
  cidr_block        = cidrsubnet(var.cidr_block, 8, 5)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-private-b"
      "Environment" = var.environment
      "Type"        = "private"
    })
  )
}

resource "aws_subnet" "private_c" {
  cidr_block        = cidrsubnet(var.cidr_block, 8, 6)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[2]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-private-c"
      "Environment" = var.environment
      "Type"        = "private"
    })
  )
}

#####################################################
# NAT Gateway (for private subnets)                #
#####################################################
resource "aws_eip" "nat" {
  count  = var.enable_nat_gateway ? 1 : 0
  domain = "vpc"

  depends_on = [aws_internet_gateway.main]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-nat-eip"
      "Environment" = var.environment
    })
  )
}

resource "aws_nat_gateway" "main" {
  count         = var.enable_nat_gateway ? 1 : 0
  allocation_id = aws_eip.nat[0].id
  subnet_id     = aws_subnet.public_a.id

  depends_on = [aws_internet_gateway.main]

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-nat"
      "Environment" = var.environment
    })
  )
}

#####################################################
# Public Route Table                               #
#####################################################
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-public"
      "Environment" = var.environment
      "Type"        = "public"
    })
  )
}

resource "aws_route" "public_internet_access" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
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

#####################################################
# Private Route Table                              #
#####################################################
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name"        = "${var.vpc_name}-${var.environment}-private"
      "Environment" = var.environment
      "Type"        = "private"
    })
  )
}

resource "aws_route" "private_nat_access" {
  count                  = var.enable_nat_gateway ? 1 : 0
  route_table_id         = aws_route_table.private.id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.main[0].id
}

# Private Route Table Associations
resource "aws_route_table_association" "private_a" {
  subnet_id      = aws_subnet.private_a.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "private_b" {
  subnet_id      = aws_subnet.private_b.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "private_c" {
  subnet_id      = aws_subnet.private_c.id
  route_table_id = aws_route_table.private.id
}
