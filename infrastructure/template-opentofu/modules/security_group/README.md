# Security Group

This module will provisioning security group for all resource used.

---

## Resources

| Name | Type |
|------|------|
| [aws_security_group.sg](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [aws_security_group_rule.egress_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |
| [aws_security_group_rule.ingress_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_sg_id"></a> [sg\_id](#output\_sg\_id) | ID of Security Group |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_common_tags"></a> [common\_tags](#input\_common\_tags) | Common tags | `map(any)` | n/a | yes |
| <a name="input_description"></a> [description](#input\_description) | Description for security group | `string` | n/a | yes |
| <a name="input_egress_rule"></a> [egress\_rule](#input\_egress\_rule) | The egress rule for security group | `any` | `[]` | no |
| <a name="input_ingress_rule"></a> [ingress\_rule](#input\_ingress\_rule) | The ingress rule for security group | `any` | `[]` | no |
| <a name="input_sg_name"></a> [sg\_name](#input\_sg\_name) | Name of security group | `string` | n/a | yes |
| <a name="input_vpc_id"></a> [vpc\_id](#input\_vpc\_id) | ID of VPC for security group | `string` | n/a | yes |

---
