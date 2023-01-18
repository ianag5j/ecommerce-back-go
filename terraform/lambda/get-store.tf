data "archive_file" "lambda_get_store_zip" {
  type        = "zip"
  source_file = "../bin/get-store/main"
  output_path = "bin/get-store/main.zip"
}

resource "aws_lambda_function" "get_store" {
  filename         = data.archive_file.lambda_get_store_zip.output_path
  function_name    = "${terraform.workspace}GetStoreGo"
  handler          = "main"
  source_code_hash = base64sha256(data.archive_file.lambda_get_store_zip.output_path)
  runtime          = "go1.x"
  role             = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/dev_serverless_lambda"
  timeout          = 15

  environment {
    variables = {
      STORES_TABLE = "${terraform.workspace}Stores"
    }
  }
}

resource "aws_lambda_permission" "api_gw_get_store" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_store.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:us-east-1:${data.aws_caller_identity.current.account_id}:${var.api_id}/*/*"
}

resource "aws_cloudwatch_log_group" "get_store" {
  name              = "/aws/lambda/${aws_lambda_function.get_store.function_name}"
  retention_in_days = 5
}

# ############ API GATEWAY ############
resource "aws_apigatewayv2_integration" "get_store" {
  api_id = var.api_id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.get_store.invoke_arn
}

resource "aws_apigatewayv2_route" "get_store" {
  api_id = var.api_id

  route_key          = "GET /v2/store"
  target             = "integrations/${aws_apigatewayv2_integration.get_store.id}"
  authorizer_id      = var.authorizer_id
  authorization_type = "CUSTOM"
}
