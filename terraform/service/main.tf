locals {
  apis = [
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "secretmanager.googleapis.com",
  ]
}

resource "google_project_service" "required" {
  for_each = toset(local.apis)

  service            = each.value
  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "app" {
  location      = var.region
  repository_id = "cep-weather"
  format        = "DOCKER"
  description   = "Images for the CEP weather API."

  depends_on = [google_project_service.required]
}

resource "google_service_account" "run" {
  account_id   = "cep-weather"
  display_name = "CEP weather Cloud Run runtime"
}

resource "google_secret_manager_secret" "weather_api_key" {
  secret_id = "cep-weather-api-key"

  replication {
    auto {}
  }

  depends_on = [google_project_service.required]
}

resource "google_secret_manager_secret_version" "weather_api_key" {
  secret      = google_secret_manager_secret.weather_api_key.id
  secret_data = var.weather_api_key
}

resource "google_secret_manager_secret_iam_member" "run" {
  secret_id = google_secret_manager_secret.weather_api_key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.run.email}"
}

resource "google_cloud_run_v2_service" "app" {
  name                = "cep-weather"
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_ALL"
  deletion_protection = false

  template {
    service_account = google_service_account.run.email
    timeout         = "30s"

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    containers {
      image = var.image

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        cpu_idle = true
      }

      env {
        name = "WEATHER_API_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.weather_api_key.secret_id
            version = "latest"
          }
        }
      }
    }
  }

  depends_on = [
    google_project_service.required,
    google_secret_manager_secret_iam_member.run,
  ]
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = google_cloud_run_v2_service.app.project
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
