# Setting Up GitHub Secrets for GCPGoLang

This document explains how to set up the required GitHub secrets for Workload Identity Federation with GCP.

## Required Secrets

You need to add the following secrets to your GitHub repository:

1. **`WIF_PROVIDER`**: The Workload Identity Provider resource name
2. **`SERVICE_ACCOUNT`**: The service account email address
3. **`GCP_PROJECT_ID`**: Your GCP project ID

## First-Time Setup

### 1. Create and Configure the Workload Identity Pool and Provider

```bash
# Set environment variables
export PROJECT_ID="gcpgolang"
export PROJECT_NUMBER="652769711122"
export GITHUB_REPO="GCPGoLang"
export GITHUB_ORG="HiepTLe"

# Create a Workload Identity Pool
gcloud iam workload-identity-pools create "github-actions-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="GitHub Actions Pool"

# Create a Workload Identity Provider in the pool
gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --display-name="GitHub Actions Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

# Get the Workload Identity Provider resource name
WORKLOAD_IDENTITY_PROVIDER=$(gcloud iam workload-identity-pools providers describe "github-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --format="value(name)")

echo "Workload Identity Provider: ${WORKLOAD_IDENTITY_PROVIDER}"
```

### 2. Create and Configure the Service Account

```bash
# Create the service account
gcloud iam service-accounts create "github-actions-sa" \
  --project="${PROJECT_ID}" \
  --display-name="GitHub Actions Service Account"

# Get the service account email
SERVICE_ACCOUNT_EMAIL="github-actions-sa@${PROJECT_ID}.iam.gserviceaccount.com"

# Grant necessary permissions to the service account
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/iam.securityReviewer"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/cloudasset.viewer"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/logging.viewer"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/storage.admin"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/pubsub.admin"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/serviceusage.serviceUsageAdmin"

# Allow GitHub to impersonate the service account
gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_EMAIL}" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository/${GITHUB_ORG}/${GITHUB_REPO}"

echo "Service Account Email: ${SERVICE_ACCOUNT_EMAIL}"
```

### 3. Add Secrets to GitHub

1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret" and add each of the following:

   - **Name**: `WIF_PROVIDER`  
     **Value**: (The Workload Identity Provider resource name from step 1)

   - **Name**: `SERVICE_ACCOUNT`  
     **Value**: (The service account email from step 2)
     
   - **Name**: `GCP_PROJECT_ID`  
     **Value**: `gcpgolang`

## Verifying the Setup

After adding the secrets, trigger the workflow manually to verify everything works:

1. Go to the Actions tab in your repository
2. Select "Terraform Infrastructure Management"
3. Click "Run workflow"
4. Choose the branch and environment
5. Click "Run workflow"

The workflow should now authenticate successfully to GCP using Workload Identity Federation.

## Bootstrapping Consideration

Once this initial setup is complete, future changes to the Workload Identity Pool, Provider, and Service Account can be managed through Terraform using the CI/CD pipeline. 