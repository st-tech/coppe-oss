resource "google_pubsub_topic" "topic_scheduler" {
  name = local.topic_name_scheduler
}

resource "google_pubsub_topic" "topic_validator" {
  name = local.topic_name_validator
}

resource "google_cloud_scheduler_job" "every_minute_job" {
  name        = "coppe-scheduler"
  description = "invoke coppe-scheduler once per minute to create coppe-apps"
  schedule    = "* * * * *"
  time_zone   = "Asia/Tokyo"

  pubsub_target {
    topic_name = google_pubsub_topic.topic_scheduler.id
    data       = base64encode("test")
  }
}

resource "google_monitoring_alert_policy" "alert_policy_scheduler_execution" {
  display_name = "Coppe-Scheduer has not been executed for last 5 min"
  combiner     = "OR"
  conditions {
    display_name = "Executions for Coppe-Scheduler [RATE]"
    condition_threshold {
      filter = "metric.type=\"cloudfunctions.googleapis.com/function/execution_count\" AND resource.type=\"cloud_function\" resource.label.\"function_name\"=\"Coppe-Scheduler\""
      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_RATE"
      }

      comparison      = "COMPARISON_LT"
      duration        = "300s"
      threshold_value = 0.01
      trigger {
        percent = 100
      }
    }
  }
  notification_channels = [var.slack_channel_id]
}

resource "google_monitoring_alert_policy" "alert_policy_error_log_from_scheduler" {
  display_name = "Detected Error Log from Coppe-Scheduler"
  combiner     = "OR"
  conditions {
    display_name = "Log entries for ERROR, Coppe-Scheduler by label.function_name [COUNT]"
    condition_threshold {
      filter = "metric.type=\"logging.googleapis.com/log_entry_count\" resource.type=\"cloud_function\" metric.label.\"severity\"=\"ERROR\" resource.label.\"function_name\"=\"Coppe-Scheduler\""
      aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_COUNT"
        group_by_fields = [
          "resource.label.function_name"
        ]
        per_series_aligner = "ALIGN_RATE"
      }

      comparison = "COMPARISON_GT"
      duration   = "0s"
      trigger {
        count = 1
      }
    }
  }
  notification_channels = [var.slack_channel_id]
}

resource "google_monitoring_alert_policy" "alert_policy_error_log_from_validator" {
  display_name = "Detected Error Log from Coppe-Validator"
  combiner     = "OR"
  conditions {
    display_name = "Log entries for ERROR, Coppe-Validator by label.function_name [COUNT]"
    condition_threshold {
      filter = "metric.type=\"logging.googleapis.com/log_entry_count\" resource.type=\"cloud_function\" metric.label.\"severity\"=\"ERROR\" resource.label.\"function_name\"=\"Coppe-Validator\""
      aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_COUNT"
        group_by_fields = [
          "resource.label.function_name"
        ]
        per_series_aligner = "ALIGN_RATE"
      }

      comparison = "COMPARISON_GT"
      duration   = "0s"
      trigger {
        count = 1
      }
    }
  }
  notification_channels = [var.slack_channel_id]
}
