data "archive_file" "lambda_authorizer_zip" {
  type        = "zip"
  source_file = "../bin/authorizer/main"
  output_path = "bin/authorizer/main.zip"
}

resource "aws_lambda_function" "authorizer" {
  filename         = data.archive_file.lambda_authorizer_zip.output_path
  function_name    = "${terraform.workspace}AuthorizerGo"
  handler          = "main"
  source_code_hash = base64sha256(data.archive_file.lambda_authorizer_zip.output_path)
  runtime          = "go1.x"
  role             = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/dev_serverless_lambda"
  timeout          = 15

  environment {
    variables = {
      AUTH0_DOMAIN = var.auth0_domain
    }
  }
}

resource "aws_cloudwatch_log_group" "authorizer" {
  name = "/aws/lambda/${aws_lambda_function.authorizer.function_name}"

  retention_in_days = 30
}

# ############ API GATEWAY ############

resource "aws_apigatewayv2_authorizer" "custom_authorizer" {
  api_id                            = var.api_id
  authorizer_payload_format_version = "2.0"
  authorizer_uri                    = aws_lambda_function.authorizer.invoke_arn
  authorizer_type                   = "REQUEST"
  identity_sources                  = ["$request.header.Authorization"]
  name                              = "ecommerce-custom-authorizer-go"
  enable_simple_responses           = true
}
