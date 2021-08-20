resource "aws_s3_bucket" "bucket" {
  for_each = toset(["bucket1", "bucket2", "bucket3"])
  bucket   = each.value
}

resource "aws_iam_user" "user" {
  for_each = toset(["user1", "user2"])
  name = each.value
}

resource "aws_iam_group" "group" {
  name = "developers"
}