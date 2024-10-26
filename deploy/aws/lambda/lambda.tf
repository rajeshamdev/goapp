
# Create Lambda function
resource "aws_lambda_function" "bowbow_lambda_func" {
  filename         = "bowbowLambdaFunc.zip"
  function_name    = "bowbow-lambda-function"
  role             = aws_iam_role.lambda_role.arn
  handler          = "main"
  source_code_hash = filebase64sha256("bowbowLambdaFunc.zip")
  runtime          = "provided.al2023"

  environment {
    variables = {
      # this variable is senstive, so it must be passed as:
      # terraform apply -var="GCP_APIKEY=<key>"
      GCP_APIKEY = var.gcp_apikey
    }
  }
}

# Create IAM role for the Lambda function
resource "aws_iam_role" "lambda_role" {

  name = "lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

# Attach a policy to the role
resource "aws_iam_role_policy_attachment" "lambda_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}
