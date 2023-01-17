data "archive_file" "lambda_create_order_zip" {
  type        = "zip"
  source_file = "../bin/create-order/main"
  output_path = "bin/create-order/main.zip"
}

resource "aws_lambda_function" "create_order" {
  filename         = data.archive_file.lambda_create_order_zip.output_path
  function_name    = "${terraform.workspace}CreateOrderGo"
  handler          = "main"
  source_code_hash = base64sha256(data.archive_file.lambda_create_order_zip.output_path)
  runtime          = "go1.x"
  role             = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/dev_serverless_lambda"
  timeout          = 15

  environment {
    variables = {
      STORES_TABLE      = "${terraform.workspace}Stores"
      ORDERS_TABLE      = "${terraform.workspace}Orders"
      CREDENTIALS_TABLE = "${terraform.workspace}Credentials"
      ENVIROMENT        = terraform.workspace
      FRONT_BASE_URL    = "https://ecommerce-front-git-development-iangonzalez-ualacomar.vercel.app"
      LAMBDA_URL        = var.base_url
    }
  }
}

resource "aws_lambda_permission" "api_gw_create_order" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.create_order.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:us-east-1:${data.aws_caller_identity.current.account_id}:${var.api_id}/*/*"
}

resource "aws_cloudwatch_log_group" "create_order" {
  name              = "/aws/lambda/${aws_lambda_function.create_order.function_name}"
  retention_in_days = 5
}

# ############ API GATEWAY ############
resource "aws_apigatewayv2_integration" "create_order" {
  api_id = var.api_id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.create_order.invoke_arn
}

resource "aws_apigatewayv2_route" "create_order" {
  api_id = var.api_id

  route_key = "POST /v2/orders"
  target    = "integrations/${aws_apigatewayv2_integration.create_order.id}"
}
