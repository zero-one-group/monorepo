# ECR

This module for stored docker image all apps

---

## Resources

| Name                                                                                                                  | Type     |
| --------------------------------------------------------------------------------------------------------------------- | -------- |
| [aws_ecr_repository.repo](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecr_repository) | resource |

## Outputs

| Name                                                                          | Description       |
| ----------------------------------------------------------------------------- | ----------------- |
| <a name="output_repository_url"></a> [repository_url](#output_repository_url) | URL of Repository |

## Inputs

| Name                                                               | Description        | Type     | Default | Required |
| ------------------------------------------------------------------ | ------------------ | -------- | ------- | :------: |
| <a name="input_common_tags"></a> [common_tags](#input_common_tags) | Common tags        | `map`    | n/a     |   yes    |
| <a name="input_repo_name"></a> [repo_name](#input_repo_name)       | Name of repository | `string` | n/a     |   yes    |

---
