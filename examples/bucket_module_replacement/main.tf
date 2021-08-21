# https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest
module "s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"
  bucket = "my-s3-bucket"
}

resource "aws_ssm_parameter" "foo" {
  name  = "foo"
  type  = "String"
  value = module.s3_bucket.s3_bucket_arn
}