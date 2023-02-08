resource "aws_dynamodb_table" "stores-dynamodb-table" {
  name         = "${terraform.workspace}Stores"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Name"

  attribute {
    name = "Name"
    type = "S"
  }

  attribute {
    name = "UserId"
    type = "S"
  }

  global_secondary_index {
    name               = "UserIdIndex"
    hash_key           = "UserId"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "INCLUDE"
    non_key_attributes = ["Name", "Id"]
  }


  tags = {
    Name        = "stores"
    Environment = terraform.workspace
  }
}
