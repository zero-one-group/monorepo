/**
 * # Keypair
 *
 * This module keypair for EC2 Instance
 */

resource "aws_key_pair" "key" {
  key_name   = var.keyname
  public_key = file(var.public_key_path)

  tags = var.common_tags
}
