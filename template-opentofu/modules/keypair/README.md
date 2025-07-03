# Keypair

This module keypair for EC2 Instance

---

## Resources

| Name                                                                                                     | Type     |
| -------------------------------------------------------------------------------------------------------- | -------- |
| [aws_key_pair.key](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/key_pair) | resource |

## Outputs

| Name                                                     | Description     |
| -------------------------------------------------------- | --------------- |
| <a name="output_keyname"></a> [keyname](#output_keyname) | Name of keypair |

## Inputs

| Name                                                                           | Description                         | Type       | Default | Required |
| ------------------------------------------------------------------------------ | ----------------------------------- | ---------- | ------- | :------: |
| <a name="input_common_tags"></a> [common_tags](#input_common_tags)             | Common tags                         | `map(any)` | n/a     |   yes    |
| <a name="input_keyname"></a> [keyname](#input_keyname)                         | Keyname of keypair                  | `string`   | n/a     |   yes    |
| <a name="input_public_key_path"></a> [public_key_path](#input_public_key_path) | Path file of public key for keypair | `string`   | n/a     |   yes    |

---
