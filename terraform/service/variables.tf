variable "project_id" {
  description = "GCP project ID."
  type        = string
}

variable "region" {
  description = "Cloud Run and Artifact Registry region."
  type        = string
  default     = "southamerica-east1"
}

variable "image" {
  description = "Container image deployed to Cloud Run, including the tag."
  type        = string
}

variable "weather_api_key" {
  description = "WeatherAPI key stored in Secret Manager."
  type        = string
  sensitive   = true
}
