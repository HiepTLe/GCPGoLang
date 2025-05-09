# Example of an insecure GCP storage bucket
# This file demonstrates configurations that would trigger policy violations

resource "google_storage_bucket" "insecure_bucket" {
  name          = "insecure-example-bucket"
  location      = "US"
  project       = var.project_id
  force_destroy = true

  # Missing uniform bucket-level access (default is false)
  # Missing versioning configuration
  # Missing encryption configuration
  # Missing logging configuration
  # Missing public access prevention

  # Public access via ACL
  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  # Overly permissive IAM policy
  # This would be defined separately with google_storage_bucket_iam_binding
}

# Separate IAM binding that makes the bucket public
resource "google_storage_bucket_iam_binding" "public_read" {
  bucket = google_storage_bucket.insecure_bucket.name
  role   = "roles/storage.objectViewer"
  members = [
    "allUsers",  # Makes the bucket publicly accessible
  ]
}

# Overly permissive IAM binding
resource "google_storage_bucket_iam_binding" "admin_access" {
  bucket = google_storage_bucket.insecure_bucket.name
  role   = "roles/storage.admin"
  members = [
    "user:example@example.com",
    "serviceAccount:service-account@${var.project_id}.iam.gserviceaccount.com",
    "group:everyone@example.com",  # Overly broad access
  ]
}

# Service account with overly permissive roles
resource "google_service_account" "insecure_sa" {
  account_id   = "insecure-service-account"
  display_name = "Insecure Service Account"
  project      = var.project_id
}

# Granting owner role to service account (violates least privilege)
resource "google_project_iam_binding" "sa_owner" {
  project = var.project_id
  role    = "roles/owner"  # Violates least privilege
  members = [
    "serviceAccount:${google_service_account.insecure_sa.email}",
  ]
} 