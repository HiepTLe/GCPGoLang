#!/bin/bash
# Script to get Workload Identity Federation information

set -e

# Set variables
PROJECT_ID="gcpgolang"
PROJECT_NUMBER="652769711122"

echo "Listing all Workload Identity Pools..."
POOLS=$(gcloud iam workload-identity-pools list --project="${PROJECT_ID}" --location="global" --format="value(name)")

if [ -z "$POOLS" ]; then
  echo "No pools found"
else
  for POOL in $POOLS; do
    POOL_NAME=$(basename "$POOL")
    echo "Found pool: $POOL_NAME"
    
    echo "Listing providers in pool $POOL_NAME..."
    PROVIDERS=$(gcloud iam workload-identity-pools providers list \
      --project="${PROJECT_ID}" \
      --location="global" \
      --workload-identity-pool="$POOL_NAME" \
      --format="value(name)" 2>/dev/null || echo "")
    
    if [ -n "$PROVIDERS" ]; then
      for PROVIDER in $PROVIDERS; do
        PROVIDER_NAME=$(basename "$PROVIDER")
        echo "Found provider: $PROVIDER_NAME"
        
        echo "Full provider resource name:"
        gcloud iam workload-identity-pools providers describe "$PROVIDER_NAME" \
          --project="${PROJECT_ID}" \
          --location="global" \
          --workload-identity-pool="$POOL_NAME" \
          --format="value(name)"
      done
    else
      echo "No providers found in pool $POOL_NAME"
    fi
  done
fi

echo "Service Account email:"
SERVICE_ACCOUNT_EMAIL="github-actions-sa@${PROJECT_ID}.iam.gserviceaccount.com"
echo "$SERVICE_ACCOUNT_EMAIL"

echo "GCP Project ID:"
echo "$PROJECT_ID" 