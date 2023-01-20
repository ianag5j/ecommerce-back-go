data "archive_file" "lambda_get_orders_zip" {
  type        = "zip"
  source_file = "../bin/get-orders/main"
  output_path = "bin/get-orders/main.zip"
}

resource "aws_lambda_function" "get_orders" {
  filename         = data.archive_file.lambda_get_orders_zip.output_path
  function_name    = "${terraform.workspace}GetOrdersGo"
  handler          = "main"
  source_code_hash = base64sha256(data.archive_file.lambda_get_orders_zip.output_path)
  runtime          = "go1.x"
  role             = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/dev_serverless_lambda"

  environment {
    variables = {
      STORES_TABLE = "${terraform.workspace}Stores"
      ORDERS_TABLE = "${terraform.workspace}Orders"
    }
  }
}

resource "aws_lambda_permission" "api_gw_get_orders" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_orders.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:us-east-1:${data.aws_caller_identity.current.account_id}:${var.api_id}/*/*"
}

resource "aws_cloudwatch_log_group" "get_orders" {
  name              = "/aws/lambda/${aws_lambda_function.get_orders.function_name}"
  retention_in_days = 5
}

# ############ API GATEWAY ############
resource "aws_apigatewayv2_integration" "get_orders" {
  api_id = var.api_id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.get_orders.invoke_arn
}

resource "aws_apigatewayv2_route" "get_orders" {
  api_id = var.api_id

  route_key          = "GET /v2/orders"
  target             = "integrations/${aws_apigatewayv2_integration.get_orders.id}"
  authorizer_id      = aws_apigatewayv2_authorizer.custom_authorizer.id
  authorization_type = "CUSTOM"
}
