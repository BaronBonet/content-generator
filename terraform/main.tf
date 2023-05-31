resource "aws_s3_bucket" "builds" {
  bucket = "${var.project_name}-builds"
}

# Required so this can be created at once, the lambda needs an image to be created so we make a dummy one for when we first apply the infra
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
  source = "terraform-aws-modules/lambda/aws"

  function_name = var.project_name
  description   = "Lambda to run the content generation pipeline"
  handler       = "index.lambda_handler"
  runtime       = "go1.x"

  create_package      = false
  s3_existing_package = {
    bucket = aws_s3_bucket.builds.id
    key    = "latest.zip"
  }
}
