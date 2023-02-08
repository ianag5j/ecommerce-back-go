resource "aws_dynamodb_table" "credentials-dynamodb-table" {
  name         = "${terraform.workspace}Credentials"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserId"
  range_key    = "Provider"

  attribute {
    name = "UserId"
    type = "S"
  }

  attribute {
    name = "Provider"
    type = "S"
  }

  tags = {
    Name        = "credentials"
    Environment = terraform.workspace
  }
}
