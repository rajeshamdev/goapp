variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "gcp_apikey" {
  description = "GCP_APIKEY"
  type        = string
  sensitive   = true
}

variable "insights_allow_cors_origins" {
  description = "list of allow cors origins"
  type = list
}

variable "insights_allow_methods" {
  description = "list of allowed methods allowed"
  type = list
}