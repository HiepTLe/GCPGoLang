/**
 * GCPGoLang GitHub Actions - Workload Identity Federation
 * 
 * This Terraform configuration sets up Workload Identity Federation between 
 * GitHub Actions and Google Cloud Platform, following security best practices.
 */

# Create a Workload Identity Pool for GitHub Actions
resource "google_iam_workload_identity_pool" "github_pool" {
  provider                  = google-beta
  project                   = var.project_id
  workload_identity_pool_id = "github-actions-pool"
  display_name              = "GitHub Actions Pool"
  description               = "Identity pool for GitHub Actions workflows"
}

# Create Workload Identity Provider for GitHub
resource "google_iam_workload_identity_pool_provider" "github_provider" {
  provider                           = google-beta
  project                            = var.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "GitHub Actions Provider"
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
    "attribute.aud"        = "assertion.aud"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

# Create Service Account for GitHub Actions
resource "google_service_account" "github_actions_sa" {
  project      = var.project_id
  account_id   = "github-actions-sa"
  display_name = "GitHub Actions Service Account"
  description  = "Service Account used by GitHub Actions workflows for GCPGoLang security tooling"
}

# Allow the GitHub repo to impersonate the service account
resource "google_service_account_iam_binding" "workload_identity_binding" {
  service_account_id = google_service_account.github_actions_sa.name
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "principalSet://iam.googleapis.com/projects/${var.project_number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.github_pool.workload_identity_pool_id}/attribute.repository/${var.github_org}/${var.github_repo}"
  ]
}

# Grant necessary permissions to the service account for IAM analysis
resource "google_project_iam_member" "iam_analyzer_permissions" {
  for_each = toset([
    "roles/iam.securityReviewer",
    "roles/cloudasset.viewer"
  ])
  
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Grant necessary permissions for Service Account Tracker
resource "google_project_iam_member" "sa_tracker_permissions" {
  for_each = toset([
    "roles/iam.securityReviewer",
    "roles/logging.viewer"
  ])
  
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Grant necessary permissions for Misconfiguration Scanner
resource "google_project_iam_member" "misconfig_scanner_permissions" {
  for_each = toset([
    "roles/cloudasset.viewer", 
    "roles/browser",
    "roles/storage.objectViewer",
    "roles/compute.viewer",
    "roles/iam.securityReviewer"
  ])
  
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Grant necessary permissions for Log Watcher
resource "google_project_iam_member" "log_watcher_permissions" {
  project = var.project_id
  role    = "roles/logging.viewer"
  member  = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Create a custom log metric for service account usage
resource "google_logging_metric" "sa_usage_metric" {
  name        = "github_actions_sa_usage"
  description = "Counts usage of the GitHub Actions service account"
  filter      = "protoPayload.authenticationInfo.principalEmail=${google_service_account.github_actions_sa.email}"
  metric_descriptor {
    metric_kind = "DELTA"
    value_type  = "INT64"
    labels {
      key         = "resource_type"
      value_type  = "STRING"
      description = "Type of resource being accessed"
    }
  }
  label_extractors = {
    "resource_type" = "EXTRACT(protoPayload.resourceName)"
  }
}

# Create an alerting policy for suspicious service account activity
resource "google_monitoring_alert_policy" "sa_alert_policy" {
  display_name = "GitHub Actions SA Suspicious Activity"
  combiner     = "OR"
  conditions {
    display_name = "High volume of operations"
    condition_threshold {
      filter          = "metric.type=\"logging.googleapis.com/user/${google_logging_metric.sa_usage_metric.name}\" resource.type=\"global\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"
      threshold_value = 100
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }

  documentation {
    content   = "There has been a high volume of activity from the GitHub Actions service account, which might indicate unauthorized use."
    mime_type = "text/markdown"
  }

  notification_channels = var.alert_notification_channels
  depends_on            = [google_logging_metric.sa_usage_metric]
}

# Output the Workload Identity Provider resource name
output "workload_identity_provider" {
  value       = google_iam_workload_identity_pool_provider.github_provider.name
  description = "The Workload Identity Provider resource name to use in GitHub Actions"
}

# Output the Service Account email
output "service_account_email" {
  value       = google_service_account.github_actions_sa.email
  description = "The Service Account email to use in GitHub Actions"
} 