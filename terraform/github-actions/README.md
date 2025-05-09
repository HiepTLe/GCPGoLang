# GCPGoLang GitHub Actions Authentication

This Terraform configuration sets up Workload Identity Federation between GitHub Actions and Google Cloud Platform for the GCPGoLang project, following enterprise security best practices.

## Overview

Workload Identity Federation allows GitHub Actions to authenticate with GCP without using static service account keys, improving security by:

- Eliminating the need to manage long-lived service account keys
- Providing fine-grained access control
- Enabling comprehensive audit logging
- Implementing proper RBAC (Role-Based Access Control)
- Setting up monitoring and alerting for suspicious activity

## Prerequisites

- A Google Cloud Platform project with billing enabled
- Terraform 1.0.0 or newer
- `gcloud` CLI configured with admin access to your project
- GitHub repository where the workflows will run

## Setup Instructions

1. **Get your GCP project information**:

   ```bash
   # Get your project ID
   PROJECT_ID=$(gcloud config get-value project)
   echo $PROJECT_ID
   
   # Get your project number
   PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
   echo $PROJECT_NUMBER
   ```

2. **Create your terraform.tfvars file**:

   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit the file with your actual values
   ```

3. **Initialize and apply the Terraform configuration**:

   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

4. **Get the output values for GitHub Actions**:

   ```bash
   terraform output
   ```

5. **Configure GitHub Repository Secrets**:

   Add the following secrets to your GitHub repository:
   
   - `WIF_PROVIDER`: Set to the value of `workload_identity_provider` from Terraform output
   - `SERVICE_ACCOUNT`: Set to the value of `service_account_email` from Terraform output
   - `GCP_PROJECT_ID`: Your Google Cloud Project ID

## Using in GitHub Actions Workflows

Update your GitHub Actions workflow files to use Workload Identity Federation:

```yaml
- name: Authenticate to Google Cloud
  id: auth
  uses: google-github-actions/auth@v2
  with:
    workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
    service_account: ${{ secrets.SERVICE_ACCOUNT }}
```

## Security Features

This configuration implements several security best practices:

1. **Least Privilege**: Service account has only the permissions it needs
2. **Short-lived Credentials**: Uses temporary OAuth tokens instead of long-lived keys
3. **Repository Binding**: Only allowed repos can use the service account
4. **Audit Logging**: All actions are logged in Cloud Audit Logs
5. **Monitoring & Alerting**: Suspicious activity triggers alerts

## Clean Up

To remove all resources created by this configuration:

```bash
terraform destroy
```

## Support

For issues or questions, please open a GitHub issue in the main GCPGoLang repository. 