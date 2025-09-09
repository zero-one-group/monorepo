/**
 * # ECR
 *
 * This module for stored docker image all apps
 */
resource "aws_ecr_repository" "repo" {
  name = var.repo_name
  tags = var.common_tags

  encryption_configuration {
    encryption_type = "KMS"
  }
  image_scanning_configuration {
    scan_on_push = true
  }
  image_tag_mutability = "IMMUTABLE"
}

# Add lifecycle policy to remove untagged images
resource "aws_ecr_lifecycle_policy" "remove_untagged" {
  repository = aws_ecr_repository.repo.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Remove untagged images"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = var.untagged_image_expiration_days
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
