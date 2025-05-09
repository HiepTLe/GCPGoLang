# Workload Identity Federation Setup

This document describes the Workload Identity Federation (WIF) setup for GitHub Actions in the GCPGoLang project.

## Overview

We use Workload Identity Federation to allow GitHub Actions to authenticate with Google Cloud Platform without using service account keys. This follows security best practices by eliminating long-lived credentials.

## Setup Options

### Option 1: Using Terraform (Recommended)

The most maintainable approach is to use our Terraform module:

```bash
cd terraform/github-actions
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your project details
terraform init
terraform apply
```

### Option 2: Using the Setup Script

For an interactive guided setup:

```bash
./scripts/setup-github-auth.sh
```

Or with custom parameters:

```bash
./scripts/setup-github-auth.sh --project-id=your-project-id --github-org=your-github-username
```

### Option 3: Manual Setup with Timestamp-Based Resources

If you need unique resource names to avoid conflicts:

```bash
./scripts/manual-setup.sh
```

## Resource Structure

The WIF setup creates these resources:

- **Pool**: A Workload Identity Pool for GitHub authentication
- **Provider**: An OIDC provider configured for GitHub Actions
- **Service Account**: A service account with necessary permissions
- **IAM Bindings**: Permissions allowing GitHub to impersonate the service account

## GitHub Secrets

The following secrets are required in the GitHub repository for workflows to authenticate with GCP:

- `WIF_PROVIDER`: The full Workload Identity Provider resource name
- `SERVICE_ACCOUNT`: The service account email
- `GCP_PROJECT_ID`: The GCP project ID

See [GitHub Secrets Setup](github-secrets-setup.md) for detailed instructions on adding these secrets.

## Troubleshooting

If you encounter issues with the WIF setup:

1. **Check current configuration**:
   ```bash
   ./scripts/get-wif-info.sh
   ```

2. **Clean up existing resources**:
   ```bash
   ./scripts/cleanup-wif.sh
   ```

3. **Create a new setup with unique names**:
   ```bash
   ./scripts/manual-setup.sh
   ```

### Common Issues

- `ALREADY_EXISTS` errors: Use the cleanup script to remove resources, then create new ones with unique names
- `NOT_FOUND` errors: There may be a delay in resource creation, wait a few minutes or check in the GCP console
- `PERMISSION_DENIED`: Ensure you're authenticated with gcloud and have the necessary permissions 