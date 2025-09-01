/**
 * # Security Group
 *
 * This module provisions multiple security groups using for_each with self-references.
 */
resource "aws_security_group" "sg" {
  for_each = var.security_groups

  name        = each.value.name
  description = each.value.description
  vpc_id      = var.vpc_id

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = each.value.name
    })
  )
}

locals {
  ingress_rules = flatten([
    for sg_key, sg_value in var.security_groups : [
      for rule_index, rule in sg_value.ingress_rule : {
        sg_key                   = sg_key
        rule_key                 = "${sg_key}-ingress-${rule_index}"
        description              = lookup(rule, "description", "Ingress Rule")
        from_port                = lookup(rule, "from_port", 0)
        to_port                  = lookup(rule, "to_port", 0)
        protocol                 = lookup(rule, "protocol", "-1")
        cidr_blocks              = lookup(rule, "cidr_blocks", null)
        source_security_group_id = lookup(rule, "source_security_group_key", null) != null ? aws_security_group.sg[lookup(rule, "source_security_group_key")].id : lookup(rule, "source_security_group_id", null)
      }
    ]
  ])

  egress_rules = flatten([
    for sg_key, sg_value in var.security_groups : [
      for rule_index, rule in sg_value.egress_rule : {
        sg_key                   = sg_key
        rule_key                 = "${sg_key}-egress-${rule_index}"
        description              = lookup(rule, "description", "Egress Rule")
        from_port                = lookup(rule, "from_port", 0)
        to_port                  = lookup(rule, "to_port", 0)
        protocol                 = lookup(rule, "protocol", "-1")
        cidr_blocks              = lookup(rule, "cidr_blocks", ["0.0.0.0/0"])
        source_security_group_id = lookup(rule, "source_security_group_key", null) != null ? aws_security_group.sg[lookup(rule, "source_security_group_key")].id : lookup(rule, "source_security_group_id", null)
      }
    ]
  ])
}

resource "aws_security_group_rule" "ingress_rule" {
  for_each = {
    for rule in local.ingress_rules : rule.rule_key => rule
  }

  type        = "ingress"
  description = each.value.description
  from_port   = each.value.from_port
  to_port     = each.value.to_port
  protocol    = each.value.protocol

  cidr_blocks = each.value.source_security_group_id == null ? each.value.cidr_blocks : null

  security_group_id        = aws_security_group.sg[each.value.sg_key].id
  source_security_group_id = each.value.source_security_group_id

  depends_on = [aws_security_group.sg]
}

resource "aws_security_group_rule" "egress_rule" {
  for_each = {
    for rule in local.egress_rules : rule.rule_key => rule
  }

  type        = "egress"
  description = each.value.description
  from_port   = each.value.from_port
  to_port     = each.value.to_port
  protocol    = each.value.protocol

  cidr_blocks = each.value.cidr_blocks

  security_group_id        = aws_security_group.sg[each.value.sg_key].id
  source_security_group_id = each.value.source_security_group_id

  depends_on = [aws_security_group.sg]
}
