terraform {
  required_version = ">= 1.9"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  backend "gcs" {
    prefix = "cep-weather/service"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
