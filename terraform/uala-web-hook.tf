data "archive_file" "lambda_uala_web_hook" {
  type        = "zip"
  source_file = "../bin/uala-web-hook/main"
  output_path = "bin/uala-web-hook/main.zip"
}

resource "aws_lambda_function" "uala_web_hook" {
  filename         = data.archive_file.lambda_uala_web_hook.output_path
  function_name    = "${terraform.workspace}UalaWebHookGo"
  handler          = "main"
  source_code_hash = base64sha256(data.archive_file.lambda_uala_web_hook.output_path)
  runtime          = "go1.x"
  role             = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/dev_serverless_lambda"

  environment {
    variables = {
      ORDERS_TABLE = "${terraform.workspace}Orders"
    }
  }
}

resource "aws_lambda_permission" "api_gw_uala_web_hook" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.uala_web_hook.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:us-east-1:${data.aws_caller_identity.current.account_id}:${data.terraform_remote_state.network.outputs.api_id}/*/*"
}

resource "aws_cloudwatch_log_group" "uala_web_hook" {
  name              = "/aws/lambda/${aws_lambda_function.uala_web_hook.function_name}"
  retention_in_days = 5
}

# ############ API GATEWAY ############
resource "aws_apigatewayv2_integration" "uala_web_hook" {
  api_id = data.terraform_remote_state.network.outputs.api_id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.uala_web_hook.invoke_arn
}

resource "aws_apigatewayv2_route" "uala_web_hook" {
  api_id = data.terraform_remote_state.network.outputs.api_id

  route_key = "POST /v2/uala-webhook/{orderId}"
  target    = "integrations/${aws_apigatewayv2_integration.uala_web_hook.id}"
}
