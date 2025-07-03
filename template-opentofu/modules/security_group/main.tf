/**
 * # Security Group
 *
 * This module will provisioning security group for all resource used.
 */

resource "aws_security_group" "sg" {
  name        = var.sg_name
  description = var.description
  vpc_id      = var.vpc_id
  tags        = var.common_tags
}

resource "aws_security_group_rule" "ingress_rule" {
  count = length(var.ingress_rule) > 0 ? length(var.ingress_rule) : 0

  type        = "ingress"
  description = lookup(var.ingress_rule[count.index], "description", "Ingress Rule")
  from_port   = lookup(var.ingress_rule[count.index], "from_port", 0)
  to_port     = lookup(var.ingress_rule[count.index], "to_port", 0)
  protocol    = lookup(var.ingress_rule[count.index], "protocol", -1)

  # Use cidr_blocks only if source_security_group_id is not specified
  cidr_blocks = lookup(var.ingress_rule[count.index], "source_security_group_id", null) == null ? (
    try(
      lookup(var.ingress_rule[count.index], "cidr_blocks", ["0.0.0.0/0"]),
      ["0.0.0.0/0"]
    )
  ) : null

  security_group_id        = aws_security_group.sg.id
  source_security_group_id = lookup(var.ingress_rule[count.index], "source_security_group_id", null)
}

resource "aws_security_group_rule" "egress_rule" {
  count = length(var.egress_rule) > 0 ? length(var.egress_rule) : 0

  type = "egress"
  description = lookup(
    var.egress_rule[count.index],
    "description",
    "Egress Rule",
  )
  from_port = lookup(
    var.egress_rule[count.index],
    "from_port",
    0,
  )
  to_port = lookup(
    var.egress_rule[count.index],
    "to_port",
    0,
  )
  protocol = lookup(
    var.egress_rule[count.index],
    "protocol",
    -1,
  )
  cidr_blocks = lookup(
    var.egress_rule[count.index],
    "cidr_blocks",
    "0.0.0.0/0",
  )
  security_group_id = aws_security_group.sg.id
  # source_security_group_id = aws_security_group.sg.id
}
