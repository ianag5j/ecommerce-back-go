terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  backend "s3" {
    bucket         = "terraform-state-ian"
    dynamodb_table = "tf-lock-table"
    key            = "eccomerce-back-go"
    region         = "us-east-1"
  }

  required_version = "~> 1.0"
}

data "terraform_remote_state" "network" {
  backend = "s3"
  config = {
    bucket         = "terraform-state-ian"
    key            = "env:/${terraform.workspace}/terraform"
    region         = "us-east-1"
    dynamodb_table = "tf-lock-table"
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_caller_identity" "current" {}

module "lambda" {
  source = "./lambda"

  api_id        = data.terraform_remote_state.network.outputs.api_id
  authorizer_id = data.terraform_remote_state.network.outputs.authorizer_id
  base_url      = data.terraform_remote_state.network.outputs.base_url
  auth0_domain  = var.auth0_domain
}

module "dynamodb" {
  source = "./dynamodb"
}
