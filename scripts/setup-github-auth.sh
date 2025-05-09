#!/bin/bash
# Setup script for GitHub Actions authentication with GCP

set -e

# Default values
DEFAULT_PROJECT_ID="gcpgolang"
DEFAULT_PROJECT_NUMBER="652769711122"
DEFAULT_GITHUB_ORG="HiepTLe"
DEFAULT_GITHUB_REPO="GCPGoLang"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --project-id)
      PROJECT_ID="$2"
      shift # past argument
      shift # past value
      ;;
    --project-number)
      PROJECT_NUMBER="$2"
      shift # past argument
      shift # past value
      ;;
    --github-org)
      GITHUB_ORG="$2"
      shift # past argument
      shift # past value
      ;;
    --github-repo)
      GITHUB_REPO="$2"
      shift # past argument
      shift # past value
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Use defaults if not provided
PROJECT_ID=${PROJECT_ID:-$DEFAULT_PROJECT_ID}
PROJECT_NUMBER=${PROJECT_NUMBER:-$DEFAULT_PROJECT_NUMBER}
GITHUB_ORG=${GITHUB_ORG:-$DEFAULT_GITHUB_ORG}
GITHUB_REPO=${GITHUB_REPO:-$DEFAULT_GITHUB_REPO}

echo "Setting up GitHub Actions authentication for GCP"
echo "- Project ID: ${PROJECT_ID}"
echo "- Project Number: ${PROJECT_NUMBER}"
echo "- GitHub Org/User: ${GITHUB_ORG}"
echo "- GitHub Repo: ${GITHUB_REPO}"
echo

# Confirm with user
read -p "Continue with these settings? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Setup cancelled"
  exit 1
fi

echo "→ Checking if gcloud is available..."
if ! command -v gcloud &> /dev/null; then
  echo "Error: gcloud CLI is not installed. Please install Google Cloud SDK."
  exit 1
fi

echo "→ Checking if logged in..."
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
  echo "You are not logged in. Please run 'gcloud auth login' first."
  exit 1
fi

echo "→ Setting project ${PROJECT_ID}..."
gcloud config set project "${PROJECT_ID}"

echo "→ Creating Workload Identity Pool..."
if ! gcloud iam workload-identity-pools describe "github-actions-pool" \
  --project="${PROJECT_ID}" --location="global" &> /dev/null; then
  gcloud iam workload-identity-pools create "github-actions-pool" \
    --project="${PROJECT_ID}" \
    --location="global" \
    --display-name="GitHub Actions Pool"
  echo "  ✓ Created pool 'github-actions-pool'"
else
  echo "  ✓ Pool 'github-actions-pool' already exists"
fi

echo "→ Creating Workload Identity Provider..."
if ! gcloud iam workload-identity-pools providers describe "github-provider" \
  --project="${PROJECT_ID}" --location="global" \
  --workload-identity-pool="github-actions-pool" &> /dev/null; then
  gcloud iam workload-identity-pools providers create-oidc "github-provider" \
    --project="${PROJECT_ID}" \
    --location="global" \
    --workload-identity-pool="github-actions-pool" \
    --display-name="GitHub Actions Provider" \
    --attribute-mapping="google.subject=assertion.sub,google.subject=assertion.repository" \
    --issuer-uri="https://token.actions.githubusercontent.com"
  echo "  ✓ Created provider 'github-provider'"
else
  echo "  ✓ Provider 'github-provider' already exists"
fi

echo "→ Getting Workload Identity Provider resource name..."
WORKLOAD_IDENTITY_PROVIDER=$(gcloud iam workload-identity-pools providers describe "github-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --format="value(name)")
echo "  ✓ Provider name: ${WORKLOAD_IDENTITY_PROVIDER}"

echo "→ Creating Service Account..."
SERVICE_ACCOUNT_EMAIL="github-actions-sa@${PROJECT_ID}.iam.gserviceaccount.com"
if ! gcloud iam service-accounts describe "${SERVICE_ACCOUNT_EMAIL}" --project="${PROJECT_ID}" &> /dev/null; then
  gcloud iam service-accounts create "github-actions-sa" \
    --project="${PROJECT_ID}" \
    --display-name="GitHub Actions Service Account"
  echo "  ✓ Created service account '${SERVICE_ACCOUNT_EMAIL}'"
else
  echo "  ✓ Service account '${SERVICE_ACCOUNT_EMAIL}' already exists"
fi

echo "→ Granting necessary permissions to service account..."
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
  echo "  - Granting ${ROLE}..."
  gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
    --role="${ROLE}" \
    --quiet
done

echo "→ Allowing GitHub to impersonate the service account..."
PRINCIPAL="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository/${GITHUB_ORG}/${GITHUB_REPO}"
gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_EMAIL}" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="${PRINCIPAL}"

echo
echo "=============================================="
echo "Setup complete! Add these secrets to GitHub:"
echo "=============================================="
echo
echo "WIF_PROVIDER:"
echo "${WORKLOAD_IDENTITY_PROVIDER}"
echo
echo "SERVICE_ACCOUNT:"
echo "${SERVICE_ACCOUNT_EMAIL}"
echo
echo "GCP_PROJECT_ID:"
echo "${PROJECT_ID}"
echo
echo "For documentation, see: docs/github-secrets-setup.md" 