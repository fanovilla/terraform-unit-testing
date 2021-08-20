# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/s3_bucket
data "aws_s3_bucket" "selected" {
  bucket = "bucket.test.com"
}

# local indirection for test overriding
locals {
  bucket = data.aws_s3_bucket.selected.bucket
}