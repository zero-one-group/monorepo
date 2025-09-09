/**
 * # Elastic IP
 *
 * This module creates shared Elastic IPs dynamically based on configuration
 */

resource "aws_eip" "this" {
  for_each = var.elastic_ips

  domain = "vpc"

  tags = merge(
    var.common_tags,
    tomap({
      "Name" = "${var.eip_name}-${each.key}"
    })
  )
}
