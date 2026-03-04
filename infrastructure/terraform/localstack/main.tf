terraform {
  required_version = ">= 1.14.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6"
    }
  }
}

variable "aws_access_key_id" {
  type = string
}

variable "aws_secret_access_key" {
  type = string
}

variable "aws_default_region" {
  type = string
}

variable "aws_s3_endpoint" {
  type = string
}

variable "aws_localstack_endpoint" {
  type = string
}

variable "aws_bucket_name_origin" {
  type = string
}

variable "aws_bucket_name_distribution" {
  type = string
}

provider "aws" {
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key
  region     = var.aws_default_region

  skip_requesting_account_id  = true
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true

  endpoints {
    s3  = var.aws_s3_endpoint
    sqs = var.aws_localstack_endpoint
  }
}

resource "aws_s3_bucket" "bucket-origin" {
  bucket = var.aws_bucket_name_origin
}

resource "aws_s3_bucket" "bucket-distribution" {
  bucket = var.aws_bucket_name_distribution
}
