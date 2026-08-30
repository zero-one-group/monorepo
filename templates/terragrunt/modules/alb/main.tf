####################
####### ALB #######
####################

# --- Application Load Balancer, one per environment.
# Creates the ALB security group, the tasks security group (allowing traffic
# from the ALB into ECS tasks), the load balancer itself, listeners and one
# target group + listener rule per service.
# The ECS module references the target group ARNs + tasks security group via
# inputs, keeping the ECS unit decoupled from load balancing. ---

resource "aws_security_group" "alb" {
  name        = "${var.project_name}-${var.environment}-alb"
  description = "Security group for the Application Load Balancer"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = var.ingress_cidr_blocks
    content {
      from_port   = var.listener_port
      to_port     = var.listener_port
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
      description = "HTTP listener ingress"
    }
  }

  dynamic "ingress" {
    for_each = var.https_enabled ? var.ingress_cidr_blocks : []
    content {
      from_port   = var.https_listener_port
      to_port     = var.https_listener_port
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
      description = "HTTPS listener ingress"
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

# One security group per ECS task, allowing traffic from the ALB on each
# service's container port. Attached to ECS services via
# tasks_security_group_id.
resource "aws_security_group" "tasks" {
  name        = "${var.project_name}-${var.environment}-ecs-alb-tasks"
  description = "Allows traffic from the ALB into ECS tasks"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = var.services
    content {
      from_port       = ingress.value.container_port
      to_port         = ingress.value.container_port
      protocol        = "tcp"
      security_groups = [aws_security_group.alb.id]
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
  name               = "${var.project_name}-${var.environment}-alb"
  internal           = var.internal
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.subnet_ids
  idle_timeout       = var.idle_timeout

  dynamic "access_logs" {
    for_each = var.access_logs_enabled ? [1] : []
    content {
      bucket  = var.access_logs_bucket
      enabled = true
    }
  }

  tags = var.common_tags
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = var.listener_port
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

resource "aws_lb_listener" "https" {
  count             = var.https_enabled ? 1 : 0
  load_balancer_arn = aws_lb.this.arn
  port              = var.https_listener_port
  protocol          = "HTTPS"
  certificate_arn   = var.certificate_arn

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
  for_each    = var.services
  name        = "${var.project_name}-${var.environment}-${each.key}"
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
  for_each     = var.services
  listener_arn = aws_lb_listener.http.arn
  priority     = index(keys(var.services), each.key) + 100

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

resource "aws_lb_listener_rule" "services_https" {
  for_each     = var.https_enabled ? var.services : {}
  listener_arn = aws_lb_listener.https[0].arn
  priority     = index(keys(var.services), each.key) + 100

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