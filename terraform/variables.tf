# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-east-1"
}

variable "auth0_domain" {
  description = "auth0 domain url."

  type    = string
  default = "https://{stage}-{auth0_id}.us.auth0.com"
}

variable "rollbar_token" {
  description = "rollbar token."

  type    = string
  default = ""
}
