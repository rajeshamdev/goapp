variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "gcp_apikey" {
  description = "GCP_APIKEY"
  type        = string
  sensitive   = true
}

variable "apigw_stage" {
  description = "API Gateway Stage"
  type        = string
}

variable "cors_origins" {
  description = "list of allow cors origins"
  type = list
}

variable "api_methods" {
  description = "list of allowed methods allowed"
  type = list
}
