resource "aws_s3_bucket" "this" {
  for_each      = var.buckets
  bucket        = each.value.bucket_name
  force_destroy = true

  tags = merge(
    {
      "Name" = each.value.bucket_name
    },
    var.common_tags
  )
}

resource "aws_s3_bucket_versioning" "this" {
  for_each = { for k, v in var.buckets : k => v if v.enable_versioning }
  bucket   = aws_s3_bucket.this[each.key].id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  for_each = var.buckets
  bucket   = aws_s3_bucket.this[each.key].id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  for_each = var.buckets
  bucket   = aws_s3_bucket.this[each.key].id

  block_public_acls       = true
  # Allow public policy if explicitly enabled OR if custom policy is provided OR if public read is enabled
  block_public_policy     = each.value.enable_public_read || each.value.allow_public_policy || each.value.custom_policy != "" ? false : true
  ignore_public_acls      = true
  # Allow public access only if public read is explicitly enabled
  restrict_public_buckets = each.value.enable_public_read ? false : true
}

resource "aws_s3_bucket_cors_configuration" "this" {
  for_each = { for k, v in var.buckets : k => v if length(v.cors_origin) > 0 }
  bucket   = aws_s3_bucket.this[each.key].id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET", "HEAD"]
    allowed_origins = each.value.cors_origin
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  for_each = { for k, v in var.buckets : k => v if v.enable_lifecycle_rule }
  bucket   = aws_s3_bucket.this[each.key].id

  rule {
    id     = "transition-to-ia"
    status = "Enabled"

    filter {
      prefix = ""
    }

    transition {
      days          = each.value.transition_days
      storage_class = "STANDARD_IA"
    }

    expiration {
      days = each.value.expiration_days
    }
  }
}

resource "aws_s3_bucket_policy" "lb_access_logs" {
  for_each = { for k, v in var.buckets : k => v if v.enable_lb_logging }
  bucket   = aws_s3_bucket.this[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "delivery.logs.amazonaws.com"
        }
        Action   = "s3:PutObject"
        Resource = "${aws_s3_bucket.this[each.key].arn}/${each.value.access_logs_prefix != "" ? each.value.access_logs_prefix : "*"}/*"
        Condition = {
          StringEquals = {
            "s3:x-amz-acl" = "bucket-owner-full-control"
          }
        }
      },
      {
        Effect = "Allow"
        Principal = {
          Service = "logdelivery.elasticloadbalancing.amazonaws.com"
        }
        Action   = "s3:PutObject"
        Resource = "${aws_s3_bucket.this[each.key].arn}/${each.value.access_logs_prefix != "" ? each.value.access_logs_prefix : "*"}/*"
        Condition = {
          StringEquals = {
            "s3:x-amz-acl" = "bucket-owner-full-control"
          }
        }
      },
      {
        Effect = "Allow"
        Principal = {
          Service = "delivery.logs.amazonaws.com"
        }
        Action   = "s3:GetBucketAcl"
        Resource = aws_s3_bucket.this[each.key].arn
      }
    ]
  })
}

resource "aws_s3_bucket_policy" "allow_access_from_bucket" {
  for_each = { for k, v in var.buckets : k => v if v.enable_public_read && v.custom_policy == "" }
  bucket   = aws_s3_bucket.this[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = "*"
        Action = [
          "s3:GetObject"
        ]
        Resource = [
          "${aws_s3_bucket.this[each.key].arn}/*"
        ]
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.this]
}

resource "aws_s3_bucket_policy" "custom_policy" {
  for_each = { for k, v in var.buckets : k => v if v.custom_policy != "" }
  bucket   = aws_s3_bucket.this[each.key].id

  policy = replace(
    each.value.custom_policy,
    "BUCKET_ARN_PLACEHOLDER",
    aws_s3_bucket.this[each.key].arn
  )

  depends_on = [aws_s3_bucket_public_access_block.this]
}
