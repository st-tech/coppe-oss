resource "google_project_iam_member" "pubsub" {
  project = var.project_id
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.service_account_for_runtime_scheduler.email}"
}

resource "google_project_iam_member" "scheduler_log_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.service_account_for_runtime_scheduler.email}"
}

resource "google_project_iam_member" "bigquery_data_viewer" {
  project = var.project_id
  role    = "roles/bigquery.dataViewer"
  member  = "serviceAccount:${google_service_account.service_account_for_runtime_validator.email}"
}

resource "google_project_iam_member" "bigquery_data_jobuser" {
  project = var.project_id
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.service_account_for_runtime_validator.email}"
}

resource "google_project_iam_member" "validator_log_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.service_account_for_runtime_validator.email}"
}
