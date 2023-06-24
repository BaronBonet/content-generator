resource "aws_s3_bucket" "builds" {
  bucket = "${var.project_name}-builds"
}

resource "null_resource" "create_temp_zip" {
  provisioner "local-exec" {
    command = "cd ${path.module}/templates/fake_zip && zip -r latest.zip ."
  }

  provisioner "local-exec" {
    command = "aws s3 cp ${path.module}/templates/fake_zip/latest.zip s3://${aws_s3_bucket.builds.bucket}/latest.zip"
  }

  depends_on = [aws_s3_bucket.builds]
}

module "lambda_function_existing_package_s3" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "4.18.0"

  function_name = var.project_name
  description   = "Lambda to run the content generation pipeline"
  handler       = "handler"
  runtime       = "provided.al2"
  architectures = ["arm64"]
  timeout       = 100

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  environment_variables = {
    NEW_YORK_TIMES_KEY          = var.env_vars.new_york_times_key
    OPENAI_KEY                  = var.env_vars.openai_key
    TWITTER_API_KEY             = var.env_vars.twitter_api_key
    TWITTER_API_KEY_SECRET      = var.env_vars.twitter_api_key_secret
    TWITTER_ACCESS_TOKEN        = var.env_vars.twitter_access_token
    TWITTER_ACCESS_TOKEN_SECRET = var.env_vars.twitter_access_token_secret
    INSTAGRAM_USERNAME          = var.env_vars.instagram_username
    INSTAGRAM_PASSWORD          = var.env_vars.instagram_password
  }

  create_package = false
  s3_existing_package = {
    bucket = aws_s3_bucket.builds.id
    key    = "latest.zip"
  }
}

resource "aws_cloudwatch_event_rule" "lambda_schedule" {
  name                = "${var.project_name}-schedule"
  description         = "Schedule for triggering lambda twice per day"
  schedule_expression = "cron(0 8,21 * * ? *)"
}

resource "aws_cloudwatch_event_target" "lambda_target" {
  rule = aws_cloudwatch_event_rule.lambda_schedule.name
  arn  = module.lambda_function_existing_package_s3.lambda_function_arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_function_existing_package_s3.lambda_function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.lambda_schedule.arn
}
