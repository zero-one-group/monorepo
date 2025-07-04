/**
 * # VPC
 *
 * This module will provision subnets for multiple environments (production, monitoring)
 * using existing VPC and Internet Gateway.
 * Each environment has 3 public and 3 private subnets across different availability zones.
 */

data "aws_vpc" "main" {
  id = "vpc-03822041f33ca8cb0"
}

data "aws_internet_gateway" "main" {
  internet_gateway_id = "igw-027095043efbe4ed1"
}

#####################################################
# Production Public Subnets                          #
#####################################################
resource "aws_subnet" "prod_public_a" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 1)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}a"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-prod-public-a"
    })
  )
}

resource "aws_subnet" "prod_public_b" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 2)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}b"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-prod-public-b"
    })
  )
}

resource "aws_subnet" "prod_public_c" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 3)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}c"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-prod-public-c"
    })
  )
}

#############################################
# Miscellanous Public Subnets                 #
#############################################
resource "aws_subnet" "misc_public_a" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 4)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}a"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-misc-public-a"
    })
  )
}

resource "aws_subnet" "misc_public_b" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 5)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}b"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-misc-public-b"
    })
  )
}

resource "aws_subnet" "misc_public_c" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 6)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}c"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-misc-public-c"
    })
  )
}

#############################################
# Development Public Subnets                #
#############################################
resource "aws_subnet" "dev_public_a" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 7)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}a"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-dev-public-a"
    })
  )
}

resource "aws_subnet" "dev_public_b" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 8)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}b"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-dev-public-b"
    })
  )
}

resource "aws_subnet" "dev_public_c" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 9)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}c"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-dev-public-c"
    })
  )
}

#############################################
# Staging Public Subnets                    #
#############################################
resource "aws_subnet" "staging_public_a" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 10)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}a"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-staging-public-a"
    })
  )
}

resource "aws_subnet" "staging_public_b" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 11)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}b"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-staging-public-b"
    })
  )
}

resource "aws_subnet" "staging_public_c" {
  cidr_block              = cidrsubnet(data.aws_vpc.main.cidr_block, 8, 12)
  map_public_ip_on_launch = true
  vpc_id                  = data.aws_vpc.main.id
  availability_zone       = "${var.region_name}c"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-staging-public-c"
    })
  )
}

# Public Route Table
resource "aws_route_table" "public" {
  vpc_id = data.aws_vpc.main.id

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-public"
    })
  )
}

resource "aws_route" "public_internet_access" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = data.aws_internet_gateway.main.id
}

# Production Public Route Table Associations
resource "aws_route_table_association" "prod_public_a" {
  subnet_id      = aws_subnet.prod_public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "prod_public_b" {
  subnet_id      = aws_subnet.prod_public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "prod_public_c" {
  subnet_id      = aws_subnet.prod_public_c.id
  route_table_id = aws_route_table.public.id
}

# Miscellanous Public Route Table Associations
resource "aws_route_table_association" "misc_public_a" {
  subnet_id      = aws_subnet.misc_public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "misc_public_b" {
  subnet_id      = aws_subnet.misc_public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "misc_public_c" {
  subnet_id      = aws_subnet.misc_public_c.id
  route_table_id = aws_route_table.public.id
}

# Development Public Route Table Associations
resource "aws_route_table_association" "dev_public_a" {
  subnet_id      = aws_subnet.dev_public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "dev_public_b" {
  subnet_id      = aws_subnet.dev_public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "dev_public_c" {
  subnet_id      = aws_subnet.dev_public_c.id
  route_table_id = aws_route_table.public.id
}

# Staging Public Route Table Associations
resource "aws_route_table_association" "staging_public_a" {
  subnet_id      = aws_subnet.staging_public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "staging_public_b" {
  subnet_id      = aws_subnet.staging_public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "staging_public_c" {
  subnet_id      = aws_subnet.staging_public_c.id
  route_table_id = aws_route_table.public.id
}
###################################################
# Elastic IP                                       #
###################################################

resource "aws_eip" "bastion" {
  domain = "vpc"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-bastion"
    })
  )
}

resource "aws_eip" "nginx" {
  domain = "vpc"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.vpc_name}-nginx"
    })
  )
}

