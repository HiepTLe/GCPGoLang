#!/bin/bash
# Setup script for configuring GitHub Actions integration with GCPGoLang

set -e

# Check if the necessary commands are available
commands=("gcloud" "jq")
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &> /dev/null; then
    echo "Error: $cmd is required but not installed."
    exit 1
  fi
done

# Set default project if none provided
if [ -z "$1" ]; then
  PROJECT_ID=$(gcloud config get-value project)
  read -p "GCP Project ID [$PROJECT_ID]: " input_project
  PROJECT_ID=${input_project:-$PROJECT_ID}
else
  PROJECT_ID=$1
fi

# Get repository info
read -p "GitHub username: " GH_USERNAME
read -p "GitHub repository name: " GH_REPO

# Confirm with the user
echo -e "\nSetting up GitHub Actions for:"
echo "- GCP Project: $PROJECT_ID"
echo "- GitHub repository: $GH_USERNAME/$GH_REPO"
read -p "Continue? (y/n): " confirm
if [ "$confirm" != "y" ]; then
  echo "Setup cancelled."
  exit 0
fi

echo -e "\n=== Creating Workload Identity Pool ==="
POOL_ID="github-actions-pool"
gcloud iam workload-identity-pools create "$POOL_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --display-name="GitHub Actions Pool" \
  --description="Identity pool for GitHub Actions" \
  2>/dev/null || echo "Pool already exists, continuing..."

echo -e "\n=== Creating Workload Identity Provider ==="
PROVIDER_ID="github-provider"
gcloud iam workload-identity-pools providers create-oidc "$PROVIDER_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="$POOL_ID" \
  --display-name="GitHub Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  2>/dev/null || echo "Provider already exists, continuing..."

echo -e "\n=== Creating Service Account ==="
SA_NAME="github-actions-sa"
SA_EMAIL="$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com"
gcloud iam service-accounts create "$SA_NAME" \
  --project="$PROJECT_ID" \
  --display-name="GitHub Actions Service Account" \
  2>/dev/null || echo "Service account already exists, continuing..."

echo -e "\n=== Granting IAM Roles ==="
# Viewer role for basic resource access
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:$SA_EMAIL" \
  --role="roles/viewer" \
  --quiet

# Security Reviewer for IAM analysis
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:$SA_EMAIL" \
  --role="roles/iam.securityReviewer" \
  --quiet

echo -e "\n=== Setting up Workload Identity Federation ==="
# Get project number
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format="value(projectNumber)")

# Allow GitHub Actions to impersonate the service account
gcloud iam service-accounts add-iam-policy-binding "$SA_EMAIL" \
  --project="$PROJECT_ID" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/$POOL_ID/attribute.repository/$GH_USERNAME/$GH_REPO" \
  --quiet

# Get the Workload Identity Provider resource name
WIF_PROVIDER=$(gcloud iam workload-identity-pools providers describe "$PROVIDER_ID" \
  --project="$PROJECT_ID" \
  --location="global" \
  --workload-identity-pool="$POOL_ID" \
  --format="value(name)")

echo -e "\n=== GitHub Actions Setup Complete ==="
echo -e "\nAdd the following secrets to your GitHub repository:"
echo "GCP_PROJECT_ID: $PROJECT_ID"
echo "WIF_PROVIDER: $WIF_PROVIDER"
echo "SERVICE_ACCOUNT: $SA_EMAIL"

echo -e "\nYou can add these secrets at:"
echo "https://github.com/$GH_USERNAME/$GH_REPO/settings/secrets/actions"
echo -e "\nFor more details, see the documentation in examples/README.md" 