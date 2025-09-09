resource "aws_iam_role" "role" {
  for_each = var.instances
  name     = "${var.project_name}-${each.key}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  tags = var.common_tags
}

# ECR Policy for each instance
resource "aws_iam_role_policy" "ecr_policy" {
  for_each = var.instances
  name     = "${var.project_name}-${each.key}-ecr-policy"
  role     = aws_iam_role.role[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetRepositoryPolicy",
          "ecr:DescribeRepositories",
          "ecr:ListImages",
          "ecr:DescribeImages",
          "ecr:BatchGetImage"
        ]
        Resource = "*"
      }
    ]
  })
}

# SSM Managed Policy attachment
resource "aws_iam_role_policy_attachment" "attach_policy_ssm" {
  for_each   = var.instances
  role       = aws_iam_role.role[each.key].name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedEC2InstanceDefaultPolicy"
}

# Dynamic S3 policy - only for instances that need S3 access
resource "aws_iam_role_policy" "s3_policy" {
  for_each = { for k, v in var.instances : k => v if length(v.s3_bucket_arns) > 0 }
  name     = "${var.project_name}-${each.key}-s3-policy"
  role     = aws_iam_role.role[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = concat(
          each.value.s3_bucket_arns,
          [for arn in each.value.s3_bucket_arns : "${arn}/*"]
        )
      }
    ]
  })
}

# EC2 describe permissions
resource "aws_iam_role_policy" "ec2_policy" {
  for_each = var.instances
  name     = "${var.project_name}-${each.key}-ec2-policy"
  role     = aws_iam_role.role[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances",
          "ec2:DescribeInstanceStatus",
          "ec2:DescribeTags",
          "ec2:DescribeVolumes",
          "ec2:DescribeVolumeStatus",
          "ec2:DescribeNetworkInterfaces",
          "ec2:DescribeAvailabilityZones",
          "ec2:DescribeRegions"
        ]
        Resource = ["*"]
      }
    ]
  })
}

resource "aws_iam_role_policy" "ssm_policy" {
  for_each = var.instances
  name     = "${var.project_name}-${each.key}-ssm-policy"
  role     = aws_iam_role.role[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = concat([
      {
        Effect = "Allow"
        Action = [
          "ssm:DescribeParameters",
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath",
          "ssm:PutParameter"
        ]
        Resource = length(each.value.ssm_parameter_paths) > 0 ? [
          for path in each.value.ssm_parameter_paths :
          "arn:aws:ssm:${var.region}:${var.aws_account_id}:parameter${path}"
        ] : ["*"]
      }
    ], length(each.value.ssm_parameter_paths) > 0 ? [
      {
        Effect = "Allow"
        Action = [
          "ssm:PutParameter"
        ]
        Resource = [
          for path in each.value.ssm_parameter_paths :
          "arn:aws:ssm:${var.region}:${var.aws_account_id}:parameter${path}"
        ]
      }
    ] : [])
  })
}

# Instance profiles
resource "aws_iam_instance_profile" "instance_profile" {
  for_each = var.instances
  name     = "${var.project_name}-${each.key}-instance-profile"
  role     = aws_iam_role.role[each.key].name
}

# EC2 Instances
resource "aws_instance" "instances" {
  for_each      = var.instances
  ami           = each.value.ami_instance
  instance_type = each.value.instance_type

  user_data = templatefile(
    "${path.module}/templates/${each.value.user_data_filename}",
    merge(each.value.user_data_vars, {
      hostname   = coalesce(each.value.hostname, each.key)
      aws_region = var.region
      project_name = var.project_name
      cluster_identifier = each.value.cluster_identifier
    })
  )

  iam_instance_profile   = aws_iam_instance_profile.instance_profile[each.key].name
  key_name               = var.keyname
  subnet_id              = each.value.subnet_id
  vpc_security_group_ids = each.value.security_group_ids

  root_block_device {
    volume_type           = each.value.volume_type
    volume_size           = each.value.volume_size
    delete_on_termination = true
    encrypted             = true
    kms_key_id            = ""
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "${var.project_name}-${each.key}"
      "Role" = each.key
    }
  )
}

resource "aws_eip_association" "eip_association" {
  for_each      = { for k, v in var.instances : k => v if v.enable_eip && v.eip_allocation_id != "" }
  instance_id   = aws_instance.instances[each.key].id
  allocation_id = each.value.eip_allocation_id
}
