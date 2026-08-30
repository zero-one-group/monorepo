/**
 * # VPC
 *
 * Creates a complete VPC for an environment: VPC, Internet Gateway,
 * 3 public subnets (for ALB/bastion), 3 private subnets (for ECS tasks)
 * and an optional NAT gateway so private resources can reach the internet.
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
      "Name" = "${var.vpc_name}-${var.environment}-vpc"
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
      "Name" = "${var.vpc_name}-${var.environment}-igw"
    })
  )
}

#####################################################
# Public Subnets                                    #
#####################################################
resource "aws_subnet" "public_a" {
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 8, var.subnet_offset + 1)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[0]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public-a"
    })
  )
}

resource "aws_subnet" "public_b" {
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 8, var.subnet_offset + 2)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[1]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public-b"
    })
  )
}

resource "aws_subnet" "public_c" {
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 8, var.subnet_offset + 3)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.main.id
  availability_zone       = data.aws_availability_zones.available.names[2]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-public-c"
    })
  )
}

#####################################################
# Private Subnets (ECS tasks, RDS, etc.)            #
#####################################################
resource "aws_subnet" "private_a" {
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, var.subnet_offset + 101)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-private-a"
    })
  )
}

resource "aws_subnet" "private_b" {
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, var.subnet_offset + 102)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-private-b"
    })
  )
}

resource "aws_subnet" "private_c" {
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, var.subnet_offset + 103)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[2]

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-private-c"
    })
  )
}

#####################################################
# Route Tables                                      #
#####################################################
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

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
  gateway_id             = aws_internet_gateway.main.id
}

resource "aws_route_table" "private" {
  count  = var.enable_nat_gateway ? 1 : 0
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-private"
    })
  )
}

#####################################################
# NAT Gateway (optional)                            #
#####################################################
resource "aws_eip" "nat" {
  count  = var.enable_nat_gateway ? 1 : 0
  domain = "vpc"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-nat"
    })
  )
}

resource "aws_nat_gateway" "main" {
  count         = var.enable_nat_gateway ? 1 : 0
  allocation_id = aws_eip.nat[0].id
  subnet_id     = aws_subnet.public_a.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-${var.environment}-nat"
    })
  )
}

resource "aws_route" "private_internet_access" {
  count                  = var.enable_nat_gateway ? 1 : 0
  route_table_id         = aws_route_table.private[0].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.main[0].id
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

# Private Route Table Associations
resource "aws_route_table_association" "private_a" {
  count          = var.enable_nat_gateway ? 1 : 0
  subnet_id      = aws_subnet.private_a.id
  route_table_id = aws_route_table.private[0].id
}

resource "aws_route_table_association" "private_b" {
  count          = var.enable_nat_gateway ? 1 : 0
  subnet_id      = aws_subnet.private_b.id
  route_table_id = aws_route_table.private[0].id
}

resource "aws_route_table_association" "private_c" {
  count          = var.enable_nat_gateway ? 1 : 0
  subnet_id      = aws_subnet.private_c.id
  route_table_id = aws_route_table.private[0].id
}
