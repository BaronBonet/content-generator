terraform {
  required_version = ">= 1.3"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2.1"
    }
  }
  backend "s3" {
    bucket = "content-generator-terraform"
    key    = "content-generator/s3/terraform.tfstate"
    region = "eu-central-1"
  }
}

provider "aws" {
  alias  = "eu-central-1"
  region = "eu-central-1"
}
