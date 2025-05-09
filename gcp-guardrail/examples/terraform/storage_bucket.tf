# Example secure GCP storage bucket
# This demonstrates best practices for creating a secure storage bucket in GCP

resource "google_storage_bucket" "secure_bucket" {
  name          = "gcp-secure-example-bucket"
  location      = "US"
  project       = var.project_id
  force_destroy = false

  # Enable uniform bucket-level access for better access control
  uniform_bucket_level_access = true

  # Enable versioning for data protection and recovery
  versioning {
    enabled = true
  }

  # Configure CMEK encryption
  encryption {
    default_kms_key_name = var.kms_key
  }

  # Enable access logs
  logging {
    log_bucket        = google_storage_bucket.logging_bucket.name
    log_object_prefix = "storage-logs"
  }

  # Enable public access prevention
  public_access_prevention = "enforced"

  # Configure lifecycle rules for cost optimization
  lifecycle_rule {
    condition {
      age = 90
    }
    action {
      type = "SetStorageClass"
      storage_class = "COLDLINE"
    }
  }

  lifecycle_rule {
    condition {
      age = 365
    }
    action {
      type = "SetStorageClass"
      storage_class = "ARCHIVE"
    }
  }

  # Retention policy for compliance requirements
  retention_policy {
    is_locked        = true
    retention_period = 2592000 # 30 days in seconds
  }

  # Labels for organization and billing tracking
  labels = {
    environment = "production"
    team        = "security"
    application = "guardrail-demo"
    owner       = "security-team"
  }
}

# Bucket for access logs
resource "google_storage_bucket" "logging_bucket" {
  name          = "gcp-secure-example-logs"
  location      = "US"
  project       = var.project_id
  force_destroy = false

  # Enable uniform bucket-level access
  uniform_bucket_level_access = true

  # Configure CMEK encryption
  encryption {
    default_kms_key_name = var.kms_key
  }

  # Enable public access prevention
  public_access_prevention = "enforced"

  # Lifecycle rule to remove old logs
  lifecycle_rule {
    condition {
      age = 180
    }
    action {
      type = "Delete"
    }
  }

  # Labels for organization and billing tracking
  labels = {
    environment = "production"
    team        = "security"
    application = "guardrail-demo-logs"
    owner       = "security-team"
  }
}

# Variables
variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "kms_key" {
  description = "The KMS key to use for bucket encryption"
  type        = string
  default     = ""
} 