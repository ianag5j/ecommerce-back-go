data "archive_file" "lambda_delete_credentials_zip" {
  type        = "zip"
  source_file = "../bin/delete-credentials/main"
  output_path = "bin/delete-credentials/main.zip"
}

resource "aws_lambda_function" "delete_credentials" {
  filename         = data.archive_file.lambda_delete_credentials_zip.output_path
  function_name    = "${terraform.workspace}DeleteCredentialsGo"
  handler          = "main"
  source_code_hash = base64sha256(data.archive_file.lambda_delete_credentials_zip.output_path)
  runtime          = "go1.x"
  role             = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/dev_serverless_lambda"

  environment {
    variables = {
      CREDENTIALS_TABLE = "${terraform.workspace}Credentials"
    }
  }
}

resource "aws_lambda_permission" "api_gw_delete_credentials" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.delete_credentials.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:us-east-1:${data.aws_caller_identity.current.account_id}:${var.api_id}/*/*"
}

resource "aws_cloudwatch_log_group" "delete_credentials" {
  name              = "/aws/lambda/${aws_lambda_function.delete_credentials.function_name}"
  retention_in_days = 5
}

# ############ API GATEWAY ############
resource "aws_apigatewayv2_integration" "delete_credentials" {
  api_id = var.api_id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.delete_credentials.invoke_arn
}

resource "aws_apigatewayv2_route" "delete_credentials" {
  api_id = var.api_id

  route_key          = "DELETE /v2/credentials"
  target             = "integrations/${aws_apigatewayv2_integration.delete_credentials.id}"
  authorizer_id      = var.authorizer_id
  authorization_type = "CUSTOM"
}
