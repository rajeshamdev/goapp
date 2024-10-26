
# Create HTTP API
resource "aws_apigatewayv2_api" "bowbow_http_api" {
  name          = "bowbow-http-api"
  protocol_type = "HTTP"
  description   = "AWS lambda for bowbow Project"

  cors_configuration {
    allow_origins  = var.cors_origins
    allow_methods  = var.api_methods
    allow_headers  = ["Content-Type", "Authorization"]
    expose_headers = ["X-Custom-Header"]
    max_age        = 3600
  }

}

# Create Lambda integration
resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.bowbow_http_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/${aws_lambda_function.bowbow_lambda_func.arn}/invocations"
  integration_method = "POST"
}

# Create routes
resource "aws_apigatewayv2_route" "video_route" {
  api_id    = aws_apigatewayv2_api.bowbow_http_api.id
  route_key = "GET /v1/api/video/{id}/insights"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "channel_route" {
  api_id    = aws_apigatewayv2_api.bowbow_http_api.id
  route_key = "GET /v1/api/channel/{id}/insights"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

# Create stage
resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.bowbow_http_api.id
  auto_deploy = true
  name        = var.apigw_stage
}

# Grant API Gateway permission to invoke the Lambda function
resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bowbow_lambda_func.function_name
  principal     = "apigateway.amazonaws.com"

  # Specify the source ARN to restrict access to a specific API Gateway
  source_arn = "${aws_apigatewayv2_api.bowbow_http_api.execution_arn}/*/*"
}

# Output the API Gateway URL
output "api_url" {
  value = "https://${aws_apigatewayv2_api.bowbow_http_api.id}.execute-api.${var.aws_region}.amazonaws.com"
}
