###################
####### ECS #######
###################

locals {
  all_entries        = merge(var.services, var.jobs)
  s3_entries         = { for k, v in local.all_entries : k => v if length(v.s3_bucket_arns) > 0 }
  ssm_entries        = { for k, v in local.all_entries : k => v if length(v.ssm_parameter_paths) > 0 }
  discovery_services = { for k, v in var.services : k => v if v.enable_service_discovery }
  dns_namespace_name = coalesce(var.dns_namespace_name, "${var.project_name}.local")
  alb_services       = { for k, v in var.services : k => v if v.enable_alb }
  alb_enabled        = length(local.alb_services) > 0
}

# --- Service discovery (Cloud Map): lets ECS tasks reach each other by a stable
# DNS name instead of a hardcoded IP, since Fargate tasks get a new private IP
# on every restart/redeploy. Only created when at least one service opts in. ---

resource "aws_service_discovery_private_dns_namespace" "this" {
  count       = length(local.discovery_services) > 0 ? 1 : 0
  name        = local.dns_namespace_name
  description = "Service discovery namespace for ${var.project_name} ECS services"
  vpc         = var.vpc_id
}

resource "aws_service_discovery_service" "this" {
  for_each = local.discovery_services
  name     = each.key

  dns_config {
    namespace_id   = aws_service_discovery_private_dns_namespace.this[0].id
    routing_policy = "MULTIVALUE"

    dns_records {
      ttl  = 10
      type = "A"
    }
  }
}

resource "aws_ecs_cluster" "this" {
  name = "${var.project_name}-ecs"
  tags = var.common_tags
}

resource "aws_ecs_cluster_capacity_providers" "this" {
  cluster_name       = aws_ecs_cluster.this.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 100
  }
}

# --- CloudWatch log groups, one per service/job entry ---

resource "aws_cloudwatch_log_group" "this" {
  for_each          = local.all_entries
  name              = "/ecs/${var.project_name}/${each.key}"
  retention_in_days = 14
  tags              = var.common_tags
}

###################
####### IAM #######
###################

# --- Execution role: shared by every task definition (ECR pull + logs).
# AmazonECSTaskExecutionRolePolicy already grants the ECR read permissions
# needed to pull images from an ECR repository. ---

resource "aws_iam_role" "execution" {
  name = "${var.project_name}-ecs-execution"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Principal = { Service = "ecs-tasks.amazonaws.com" }
      }
    ]
  })

  tags = var.common_tags
}

resource "aws_iam_role_policy_attachment" "execution_managed" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# --- Task roles: one per service/job entry ---

resource "aws_iam_role" "services_task" {
  for_each = var.services
  name     = "${var.project_name}-ecs-svc-${each.key}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Principal = { Service = "ecs-tasks.amazonaws.com" }
      }
    ]
  })

  tags = var.common_tags
}

resource "aws_iam_role" "jobs_task" {
  for_each = var.jobs
  name     = "${var.project_name}-ecs-job-${each.key}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Principal = { Service = "ecs-tasks.amazonaws.com" }
      }
    ]
  })

  tags = var.common_tags
}

locals {
  task_role_ids = merge(
    { for k, v in aws_iam_role.services_task : k => v.id },
    { for k, v in aws_iam_role.jobs_task : k => v.id },
  )
}

resource "aws_iam_role_policy" "s3_policy" {
  for_each = local.s3_entries
  name     = "${var.project_name}-ecs-${each.key}-s3-policy"
  role     = local.task_role_ids[each.key]

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

resource "aws_iam_role_policy" "ssm_policy" {
  for_each = local.ssm_entries
  name     = "${var.project_name}-ecs-${each.key}-ssm-policy"
  role     = local.task_role_ids[each.key]

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:DescribeParameters",
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath",
          "ssm:PutParameter"
        ]
        Resource = [
          for path in each.value.ssm_parameter_paths :
          "arn:aws:ssm:${var.region}:${var.aws_account_id}:parameter${path}"
        ]
      }
    ]
  })
}

#############################
### TASK DEFINITIONS #######
#############################

resource "aws_ecs_task_definition" "services" {
  for_each                 = var.services
  family                   = "${var.project_name}-${each.key}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = tostring(each.value.cpu)
  memory                   = tostring(each.value.memory)
  execution_role_arn       = aws_iam_role.execution.arn
  task_role_arn            = aws_iam_role.services_task[each.key].arn

  container_definitions = jsonencode([
    merge(
      {
        name        = each.key
        image       = coalesce(each.value.image, "${var.ecr_repository_url}:${var.image_tag}-${each.key}")
        essential   = true
        environment = [for k, v in each.value.environment : { name = k, value = v }]
        logConfiguration = {
          logDriver = "awslogs"
          options = {
            "awslogs-group"         = aws_cloudwatch_log_group.this[each.key].name
            "awslogs-region"        = var.region
            "awslogs-stream-prefix" = each.key
          }
        }
      },
      length(each.value.command) > 0 ? { command = each.value.command } : {}
    )
  ])

  tags = var.common_tags
}

resource "aws_ecs_task_definition" "jobs" {
  for_each                 = var.jobs
  family                   = "${var.project_name}-${each.key}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = tostring(each.value.cpu)
  memory                   = tostring(each.value.memory)
  execution_role_arn       = aws_iam_role.execution.arn
  task_role_arn            = aws_iam_role.jobs_task[each.key].arn

  container_definitions = jsonencode([
    merge(
      {
        name        = each.key
        image       = coalesce(each.value.image, "${var.ecr_repository_url}:${var.image_tag}-${each.key}")
        essential   = true
        environment = [for k, v in each.value.environment : { name = k, value = v }]
        logConfiguration = {
          logDriver = "awslogs"
          options = {
            "awslogs-group"         = aws_cloudwatch_log_group.this[each.key].name
            "awslogs-region"        = var.region
            "awslogs-stream-prefix" = each.key
          }
        }
      },
      length(each.value.command) > 0 ? { command = each.value.command } : {}
    )
  ])

  tags = var.common_tags
}

####################
### ECS SERVICES ###
####################

resource "aws_ecs_service" "services" {
  for_each        = var.services
  name            = "${var.project_name}-${each.key}"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.services[each.key].arn
  desired_count   = each.value.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.subnet_ids
    security_groups  = concat(var.security_group_ids, each.value.enable_alb ? [aws_security_group.tasks_alb[0].id] : [])
    assign_public_ip = each.value.assign_public_ip
  }

  dynamic "service_registries" {
    for_each = each.value.enable_service_discovery ? [1] : []
    content {
      registry_arn = aws_service_discovery_service.this[each.key].arn
    }
  }

  dynamic "load_balancer" {
    for_each = each.value.enable_alb ? [1] : []
    content {
      target_group_arn = aws_lb_target_group.services[each.key].arn
      container_name   = each.key
      container_port   = each.value.container_port
    }
  }

  tags = var.common_tags
}

################
####### ALB ########
################

# --- Optional Application Load Balancer, created once per module instance.
# Any service with enable_alb = true gets a target group + listener rule.
# The module owns the ALB + task security groups so there is no circular
# dependency with the external security-group module. ---

resource "aws_security_group" "alb" {
  count       = local.alb_enabled ? 1 : 0
  name        = "${var.project_name}-ecs-alb"
  description = "Security group for the ECS Application Load Balancer"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = var.alb_ingress_cidr_blocks
    content {
      from_port   = var.alb_listener_port
      to_port     = var.alb_listener_port
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
      description = "ALB listener ingress"
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = var.common_tags
}

resource "aws_security_group" "tasks_alb" {
  count       = local.alb_enabled ? 1 : 0
  name        = "${var.project_name}-ecs-alb-tasks"
  description = "Allows traffic from the ECS Application Load Balancer into ECS tasks"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = local.alb_services
    content {
      from_port       = ingress.value.container_port
      to_port         = ingress.value.container_port
      protocol        = "tcp"
      security_groups = [aws_security_group.alb[0].id]
      description     = "Traffic from the ALB on container port ${ingress.value.container_port}"
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = var.common_tags
}

resource "aws_lb" "this" {
  count              = local.alb_enabled ? 1 : 0
  name               = "${var.project_name}-ecs-alb"
  internal           = var.alb_internal
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb[0].id]
  subnets            = var.alb_subnet_ids
  idle_timeout       = var.alb_idle_timeout
  tags               = var.common_tags
}

resource "aws_lb_listener" "this" {
  count             = local.alb_enabled ? 1 : 0
  load_balancer_arn = aws_lb.this[0].arn
  port              = var.alb_listener_port
  protocol          = "HTTP"

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = "404"
    }
  }
}

resource "aws_lb_target_group" "services" {
  for_each    = local.alb_services
  name        = "${var.project_name}-${each.key}"
  port        = each.value.container_port
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    path                = each.value.alb_health_check_path
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200-399"
  }

  tags = var.common_tags
}

resource "aws_lb_listener_rule" "services" {
  for_each     = local.alb_services
  listener_arn = aws_lb_listener.this[0].arn
  priority     = index(keys(local.alb_services), each.key) + 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.services[each.key].arn
  }

  condition {
    path_pattern {
      values = [each.value.alb_path_pattern]
    }
  }
}
