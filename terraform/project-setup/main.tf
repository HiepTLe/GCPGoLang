/**
 * GCPGoLang Project Setup
 * 
 * This Terraform configuration handles the basic setup of a GCP project for GCPGoLang,
 * including enabling necessary APIs and setting up initial resources.
 */

locals {
  required_apis = [
    "iam.googleapis.com",             # Identity and Access Management API
    "cloudasset.googleapis.com",      # Cloud Asset API
    "logging.googleapis.com",         # Cloud Logging API
    "monitoring.googleapis.com",      # Cloud Monitoring API
    "pubsub.googleapis.com",          # Pub/Sub API
    "container.googleapis.com",       # Kubernetes Engine API
    "compute.googleapis.com",         # Compute Engine API
    "storage.googleapis.com",         # Cloud Storage API
    "cloudresourcemanager.googleapis.com", # Cloud Resource Manager API
    "serviceusage.googleapis.com",    # Service Usage API
    "bigquery.googleapis.com",        # BigQuery API
    "cloudkms.googleapis.com",        # Cloud KMS API
    "datastore.googleapis.com",       # Datastore API
    "secretmanager.googleapis.com",   # Secret Manager API
  ]
}

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset(local.required_apis)
  
  project = var.project_id
  service = each.value
  
  disable_dependent_services = false
  disable_on_destroy         = false
}

# Create a custom audit log sink for security monitoring
resource "google_logging_project_sink" "security_audit_sink" {
  name        = "security-audit-sink"
  description = "Security audit log sink for GCPGoLang"
  
  # This sink will export logs to a Cloud Storage bucket
  destination = "storage.googleapis.com/${google_storage_bucket.audit_logs.name}"
  
  # Use a filter that captures security-relevant logs
  filter = "logName:\"projects/${var.project_id}/logs/cloudaudit.googleapis.com%2Factivity\" OR logName:\"projects/${var.project_id}/logs/cloudaudit.googleapis.com%2Fdata_access\""
  
  # Unique writer identity for this sink
  unique_writer_identity = true
  
  depends_on = [
    google_project_service.required_apis,
    google_storage_bucket.audit_logs
  ]
}

# Create a storage bucket for audit logs
resource "google_storage_bucket" "audit_logs" {
  name     = "${var.project_id}-audit-logs"
  location = var.region
  
  # Force destroy should be false in production
  force_destroy = var.environment == "prod" ? false : true
  
  # Uniform bucket-level access is required for security
  uniform_bucket_level_access = true
  
  # Set appropriate storage class
  storage_class = "STANDARD"
  
  # Versioning to prevent data loss
  versioning {
    enabled = true
  }
  
  # Lifecycle policy for audit logs
  lifecycle_rule {
    condition {
      age = var.environment == "prod" ? 365 : 30 # Retain logs for 1 year in prod
    }
    action {
      type = "Delete"
    }
  }
  
  depends_on = [google_project_service.required_apis]
}

# IAM binding for the audit logs bucket
resource "google_storage_bucket_iam_binding" "audit_logs_writer" {
  bucket = google_storage_bucket.audit_logs.name
  role   = "roles/storage.objectCreator"
  
  members = [
    google_logging_project_sink.security_audit_sink.writer_identity,
  ]
}

# Create a Pub/Sub topic for security alerts
resource "google_pubsub_topic" "security_alerts" {
  name = "security-alerts"
  
  # Enable message retention
  message_retention_duration = "604800s" # 7 days
  
  depends_on = [google_project_service.required_apis]
}

# Output the enabled APIs
output "enabled_apis" {
  value       = [for api in google_project_service.required_apis : api.service]
  description = "List of APIs enabled for the project"
}

# Output the audit logs bucket
output "audit_logs_bucket" {
  value       = google_storage_bucket.audit_logs.name
  description = "Bucket storing security audit logs"
}

# Output the security alerts topic
output "security_alerts_topic" {
  value       = google_pubsub_topic.security_alerts.name
  description = "Pub/Sub topic for security alerts"
} 