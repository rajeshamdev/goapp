variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "gcp_apikey" {
  description = "GCP API KEY"
  type        = string
  sensitive   = true
}
