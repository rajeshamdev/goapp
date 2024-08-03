variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "gcp_apikey" {
  description = "GCP_APIKEY"
  type        = string
  sensitive   = true
}
