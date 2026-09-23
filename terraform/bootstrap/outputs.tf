output "state_bucket" {
  description = "GCS bucket used as the service stack backend."
  value       = google_storage_bucket.tfstate.name
}

output "workload_identity_provider" {
  description = "Value for the GCP_WORKLOAD_IDENTITY_PROVIDER GitHub secret."
  value       = google_iam_workload_identity_pool_provider.github.name
}

output "deployer_service_account" {
  description = "Value for the GCP_SERVICE_ACCOUNT GitHub secret."
  value       = google_service_account.deployer.email
}
