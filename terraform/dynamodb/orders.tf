resource "aws_dynamodb_table" "orders-dynamodb-table" {
  name         = "${terraform.workspace}Orders"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "StoreId"
    type = "S"
  }

  attribute {
    name = "ExternalId"
    type = "S"
  }

  global_secondary_index {
    name            = "StoreIdIndex"
    hash_key        = "StoreId"
    write_capacity  = 10
    read_capacity   = 10
    projection_type = "ALL"
  }

  global_secondary_index {
    name               = "ExternalIdIndex"
    hash_key           = "ExternalId"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "INCLUDE"
    non_key_attributes = ["Id"]
  }


  tags = {
    Name        = "orders"
    Environment = terraform.workspace
  }
}
