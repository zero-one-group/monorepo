resource "aws_iam_role" "role" {
  name               = var.instance_name
  assume_role_policy = file("${path.module}/../templates/ec2/ec2-instance-profile-policy.json")

  tags = var.common_tags
}

resource "aws_iam_role_policy" "ecr_policy" {
  name = "${var.instance_name}-ecr-policy"
  role = aws_iam_role.role.id

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

resource "aws_iam_role_policy_attachment" "attach_policy_ssm" {
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedEC2InstanceDefaultPolicy"
}

# Adjust on lgtm part based on your project name
resource "aws_iam_role_policy" "s3_lgtm_policy" {
  name = "${var.instance_name}-s3-lgtm-policy"
  role = aws_iam_role.role.id

  policy = jsonencode({
    Version = "2012-10-17"
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "s3:ListBucket",
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        "Resource" : [
          "arn:aws:s3:::lgtm-buckets-monitoring",
          "arn:aws:s3:::lgtm-buckets-monitoring/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy" "ec2_policy" {
  name = "${var.instance_name}-ec2-policy"
  role = aws_iam_role.role.id

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
      },
      {
        Sid    = "AllowInstanceAndLogDescriptions"
        Effect = "Allow"
        Action = [
          "rds:DescribeDBInstances",
          "rds:DescribeDBLogFiles"
        ]
        Resource = ["arn:aws:rds:*:*:db:*"]
      },
      {
        Sid    = "AllowMaintenanceDescriptions"
        Effect = "Allow"
        Action = [
          "rds:DescribePendingMaintenanceActions"
        ]
        Resource = "*"
      },
      {
        Sid    = "AllowGettingCloudWatchMetrics"
        Effect = "Allow"
        Action = [
          "cloudwatch:GetMetricData"
        ]
        Resource = "*"
      },
      {
        Sid    = "AllowRDSUsageDescriptions"
        Effect = "Allow"
        Action = [
          "rds:DescribeAccountAttributes"
        ]
        Resource = "*"
      },
      {
        Sid    = "AllowQuotaDescriptions"
        Effect = "Allow"
        Action = [
          "servicequotas:GetServiceQuota"
        ]
        Resource = "*"
      },
      {
        Sid    = "AllowInstanceTypeDescriptions"
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstanceTypes"
        ]
        Resource = "*"
      }
    ]
  })
}

# Adjust on lgtm part based on your project name
resource "aws_iam_role_policy" "ssm_policy" {
  name = "${var.instance_name}-ssm-policy"
  role = aws_iam_role.role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:PutParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath"
        ]
        Resource = [
          "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter/lgtm/swarm/*"
        ]
      }
    ]
  })
}

resource "aws_iam_instance_profile" "instance_profile" {
  name = "${var.instance_name}-instance-profile"
  role = aws_iam_role.role.name
}

resource "aws_instance" "instance" {
  ami           = var.ami_instance
  instance_type = var.instance_type

  user_data = templatefile(
    "${path.module}/../templates/ec2/${var.user_data.filename}",
    merge(var.user_data.vars, {
      hostname = coalesce(var.hostname, join("-", slice(split("-", var.instance_name), 1, length(split("-", var.instance_name)))))
    })
  )
  iam_instance_profile   = aws_iam_instance_profile.instance_profile.name
  key_name               = var.keyname
  subnet_id              = var.subnet
  vpc_security_group_ids = var.security_group

  root_block_device {
    volume_type           = var.root_block_device.volume_type
    volume_size           = var.root_block_device.volume_size
    delete_on_termination = var.root_block_device.delete_on_termination
    encrypted             = var.root_block_device.encrypted
    kms_key_id            = var.root_block_device.kms_key_id
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.instance_name}"
    })
  )
}
