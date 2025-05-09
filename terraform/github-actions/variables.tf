/**
 * GCPGoLang GitHub Actions - Variables
 * 
 * Input variables for the Workload Identity Federation Terraform configuration.
 */

variable "project_id" {
  description = "The GCP project ID where resources will be created"
  type        = string
}

variable "project_number" {
  description = "The GCP project number"
  type        = string
}

variable "github_org" {
  description = "The GitHub organization or username that owns the repository"
  type        = string
}

variable "github_repo" {
  description = "The GitHub repository name"
  type        = string
}

variable "alert_notification_channels" {
  description = "List of notification channel IDs for alerting"
  type        = list(string)
  default     = []
} 