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

  required_version = "~> 1.0"

  backend "s3" {
    bucket = "terraform-state-ian"
    key    = "terraform"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}
