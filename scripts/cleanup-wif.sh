#!/bin/bash
# Cleanup script for Workload Identity Federation resources

set -e

# Set variables
PROJECT_ID="gcpgolang"
PROJECT_NUMBER="652769711122"
SERVICE_ACCOUNT_NAME="github-actions-sa"
SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo "Listing existing pools..."
POOLS=$(gcloud iam workload-identity-pools list --project="${PROJECT_ID}" --location="global" --format="value(name)")

if [ -z "$POOLS" ]; then
  echo "No pools found"
else
  for POOL in $POOLS; do
    POOL_NAME=$(basename "$POOL")
    echo "Found pool: $POOL_NAME"
    
    # List and delete providers in the pool
    PROVIDERS=$(gcloud iam workload-identity-pools providers list \
      --project="${PROJECT_ID}" \
      --location="global" \
      --workload-identity-pool="$POOL_NAME" \
      --format="value(name)" 2>/dev/null || echo "")
    
    if [ -n "$PROVIDERS" ]; then
      for PROVIDER in $PROVIDERS; do
        PROVIDER_NAME=$(basename "$PROVIDER")
        echo "Deleting provider: $PROVIDER_NAME in pool $POOL_NAME"
        gcloud iam workload-identity-pools providers delete "$PROVIDER_NAME" \
          --project="${PROJECT_ID}" \
          --location="global" \
          --workload-identity-pool="$POOL_NAME" \
          --quiet
      done
    fi
    
    echo "Deleting pool: $POOL_NAME"
    gcloud iam workload-identity-pools delete "$POOL_NAME" \
      --project="${PROJECT_ID}" \
      --location="global" \
      --quiet
  done
fi

# Try to explicitly check for github-pool
echo "Attempting to force delete previously created pools..."
gcloud iam workload-identity-pools delete "github-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --quiet 2>/dev/null || echo "github-pool not found or already deleted"

echo "Cleanup completed successfully." 