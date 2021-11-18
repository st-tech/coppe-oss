resource "google_service_account" "service_account_for_runtime_scheduler" {
  account_id   = "coppe-scheduler-runtime"
  display_name = "For publishing Topic message from Coppe-Scheduler"
}

resource "google_service_account" "service_account_for_runtime_validator" {
  account_id   = "coppe-validator-runtime"
  display_name = "For accessing BigQuery from Coppe-Validator"
}
