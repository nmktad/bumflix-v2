terraform {
  required_version = ">= 1.14.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6"
    }
  }
}

variable "access_key_id" {
  type = string
}

variable "secret_access_key" {
  type = string
}

variable "default_region" {
  type = string
}

variable "localstack_s3_endpoint" {
  type = string
}

variable "localstack_endpoint" {
  type = string
}

variable "bucket_name_origin" {
  type = string
}

variable "bucket_name_distribution" {
  type = string
}

variable "ingestion_queue" {
  type = string
}

variable "deadletter_queue" {
  type = string
}

provider "aws" {
  access_key = var.access_key_id
  secret_key = var.secret_access_key
  region     = var.default_region

  skip_requesting_account_id  = true
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true

  endpoints {
    s3  = var.localstack_s3_endpoint
    sqs = var.localstack_endpoint
  }
}

resource "aws_s3_bucket" "bucket-origin" {
  bucket = var.bucket_name_origin
}

resource "aws_s3_bucket" "bucket-distribution" {
  bucket = var.bucket_name_distribution
}

resource "aws_sqs_queue" "terraform_queue" {
  name                      = var.ingestion_queue
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.terraform_queue_deadletter.arn
    maxReceiveCount     = 4
  })
}

resource "aws_sqs_queue" "terraform_queue_deadletter" {
  name = var.deadletter_queue
}

resource "aws_sqs_queue_redrive_allow_policy" "terraform_queue_redrive_allow_policy" {
  queue_url = aws_sqs_queue.terraform_queue_deadletter.id

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.terraform_queue.arn]
  })
}
