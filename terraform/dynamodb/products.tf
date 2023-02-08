resource "aws_dynamodb_table" "products-dynamodb-table" {
  name         = "${terraform.workspace}Products"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Id"

  attribute {
    name = "Id"
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
    non_key_attributes = ["Id", "Name", "Price"]
  }


  tags = {
    Name        = "products"
    Environment = terraform.workspace
  }
}
