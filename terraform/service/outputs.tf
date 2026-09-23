output "service_url" {
  description = "Public Cloud Run URL. Copy this into the README after the first deploy."
  value       = google_cloud_run_v2_service.app.uri
}

output "image_repository" {
  description = "Artifact Registry repository path without a tag."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/cep-weather/api"
}
