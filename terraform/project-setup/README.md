# GCPGoLang Project Setup

This Terraform module handles the basic setup of a GCP project for the GCPGoLang security governance toolkit, following enterprise best practices.

## Features

This module:

1. **Enables Required APIs** - Activates all necessary Google Cloud APIs
2. **Sets Up Audit Logging** - Creates a log sink to capture security-relevant events
3. **Creates a Pub/Sub Topic** - For security alerts and notifications
4. **Implements Secure Defaults** - Following Google Cloud security best practices

## Prerequisites

- A Google Cloud Platform project
- Terraform 1.0.0 or newer
- `gcloud` CLI configured with admin access to your project

## Usage

1. **Initialize your configuration**:

   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your project details
   ```

2. **Apply the Terraform configuration**:

   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Security Considerations

This module implements several security best practices:

- **Comprehensive API Enablement**: Only necessary APIs are enabled
- **Audit Logging**: Security-relevant logs are preserved
- **Secure Defaults**: All resources use secure configurations
- **Principle of Least Privilege**: IAM bindings follow least privilege

## Integration with GCPGoLang

After applying this project setup, you'll have the foundation to:

1. Deploy the GCPGoLang security tooling
2. Set up Workload Identity Federation for CI/CD
3. Begin monitoring your GCP environment for security issues

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| project_id | GCP Project ID | string | - | yes |
| project_number | GCP Project Number | string | - | yes |
| region | Default region for resources | string | us-central1 | no |
| zone | Default zone for resources | string | us-central1-a | no |
| environment | Environment type (dev, staging, prod) | string | dev | no |

## Outputs

| Name | Description |
|------|-------------|
| enabled_apis | List of enabled APIs |
| audit_logs_bucket | Bucket for security audit logs |
| security_alerts_topic | Pub/Sub topic for security alerts | 