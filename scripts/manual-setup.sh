#!/bin/bash
# Manual setup script for GitHub Actions WIF

set -e

# Set variables
PROJECT_ID="gcpgolang"
PROJECT_NUMBER="652769711122"
GITHUB_ORG="HiepTLe"
GITHUB_REPO="GCPGoLang"
TIMESTAMP=$(date +%m%d%H%M)
POOL_NAME="github-pool-${TIMESTAMP}"
PROVIDER_NAME="github-provider-${TIMESTAMP}"

echo "Creating Workload Identity Pool..."
gcloud iam workload-identity-pools create "${POOL_NAME}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="GH Pool ${TIMESTAMP}"

echo "Creating Workload Identity Provider..."
gcloud iam workload-identity-pools providers create-oidc "${PROVIDER_NAME}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="${POOL_NAME}" \
  --display-name="GH Provider ${TIMESTAMP}" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-condition="attribute.repository==\"${GITHUB_ORG}/${GITHUB_REPO}\""

echo "Getting Workload Identity Provider resource name..."
WORKLOAD_IDENTITY_PROVIDER=$(gcloud iam workload-identity-pools providers describe "${PROVIDER_NAME}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="${POOL_NAME}" \
  --format="value(name)")

echo "Creating Service Account..."
SERVICE_ACCOUNT_NAME="github-actions-sa"
SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# Check if service account exists
if gcloud iam service-accounts describe "${SERVICE_ACCOUNT_EMAIL}" --project="${PROJECT_ID}" &>/dev/null; then
  echo "Service account ${SERVICE_ACCOUNT_EMAIL} already exists, skipping creation"
else
  gcloud iam service-accounts create "${SERVICE_ACCOUNT_NAME}" \
    --project="${PROJECT_ID}" \
    --display-name="GitHub Actions Service Account"
fi

echo "Granting necessary permissions..."
# List of roles needed
ROLES=(
  "roles/iam.securityReviewer"
  "roles/cloudasset.viewer"
  "roles/logging.viewer"
  "roles/storage.admin"
  "roles/pubsub.admin"
  "roles/serviceusage.serviceUsageAdmin"
)

for ROLE in "${ROLES[@]}"; do
  echo "Granting ${ROLE}..."
  gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
    --role="${ROLE}"
done

echo "Adding workload identity binding..."
gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_EMAIL}" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_NAME}/attribute.repository/${GITHUB_ORG}/${GITHUB_REPO}"

echo "===== GitHub Secrets to Add ====="
echo "WIF_PROVIDER:"
echo "${WORKLOAD_IDENTITY_PROVIDER}"
echo
echo "SERVICE_ACCOUNT:"
echo "${SERVICE_ACCOUNT_EMAIL}"
echo
echo "GCP_PROJECT_ID:"
echo "${PROJECT_ID}" 