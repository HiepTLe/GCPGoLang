# Workload Identity Federation Setup

This document describes the Workload Identity Federation (WIF) setup for GitHub Actions in the GCPGoLang project.

## Overview

We use Workload Identity Federation to allow GitHub Actions to authenticate with Google Cloud Platform without using service account keys. This follows security best practices by eliminating long-lived credentials.

## Current Configuration

The current Workload Identity Federation setup uses:

- **Pool**: `github-pool-05091617` 
- **Provider**: `github-provider-05091617`
- **Resource Name**: `projects/652769711122/locations/global/workloadIdentityPools/github-pool-05091617/providers/github-provider-05091617`
- **Service Account**: `github-actions-sa@gcpgolang.iam.gserviceaccount.com`
- **Project ID**: `gcpgolang`
- **Project Number**: `652769711122`

## GitHub Secrets

The following secrets are required in the GitHub repository for workflows to authenticate with GCP:

- `WIF_PROVIDER`: The full Workload Identity Provider resource name
- `SERVICE_ACCOUNT`: The service account email
- `GCP_PROJECT_ID`: The GCP project ID

## Troubleshooting

If you encounter issues with the WIF setup:

1. Run `./scripts/get-wif-info.sh` to view current configuration
2. Run `./scripts/cleanup-wif.sh` to remove existing resources
3. Run `./scripts/manual-setup.sh` to recreate resources with unique names
4. Update GitHub secrets with the new values

### Common Issues

- `ALREADY_EXISTS` errors: Use the cleanup script to remove resources, then create new ones with unique names
- `NOT_FOUND` errors: There may be a delay in resource creation, wait a few minutes or check in the GCP console 