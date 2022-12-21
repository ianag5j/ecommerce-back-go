# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-east-1"
}

variable "api_id" {
  description = "aws api id."

  type = string
}

variable "authorizer_id" {
  description = "aws lambda auth id."

  type = string
}
