data "aws_caller_identity" "current" {}

data "terraform_remote_state" "network" {
  backend = "s3"
  config = {
    bucket         = "terraform-state-ian"
    key            = "env:/${terraform.workspace}/terraform"
    region         = "us-east-1"
    dynamodb_table = "tf-lock-table"
  }
}
