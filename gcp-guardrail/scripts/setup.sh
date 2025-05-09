#!/bin/bash
# Setup script for GCP GuardRail

set -e

echo "Setting up GCP GuardRail..."

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "Go is not installed. Please install Go 1.21+"; exit 1; }
command -v gcloud >/dev/null 2>&1 || { echo "gcloud is not installed. Please install Google Cloud SDK"; exit 1; }

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
    echo "Go version 1.21+ is required. Current version: $GO_VERSION"
    exit 1
fi

# Build the tools
echo "Building tools..."
mkdir -p bin
go build -o bin/iam-analyzer ./cmd/iam-analyzer
go build -o bin/tf-validator ./cmd/tf-validator
go build -o bin/log-watcher ./cmd/log-watcher
go build -o bin/sa-tracker ./cmd/sa-tracker

echo "Setting up playground..."
mkdir -p playground/static
mkdir -p playground/templates
go build -o bin/rego-playground ./playground

echo "Setting up controller..."
go build -o bin/admission-controller ./controller

# Setup environment
echo "Setting up GCP environment..."
read -p "Enter GCP project ID: " PROJECT_ID

# Check if the project exists and if the user has access
gcloud projects describe $PROJECT_ID >/dev/null 2>&1 || { echo "Project $PROJECT_ID not found or you don't have access"; exit 1; }

# Enable required APIs
echo "Enabling required GCP APIs..."
gcloud services enable iam.googleapis.com --project $PROJECT_ID
gcloud services enable cloudasset.googleapis.com --project $PROJECT_ID
gcloud services enable logging.googleapis.com --project $PROJECT_ID
gcloud services enable monitoring.googleapis.com --project $PROJECT_ID
gcloud services enable pubsub.googleapis.com --project $PROJECT_ID
gcloud services enable container.googleapis.com --project $PROJECT_ID
gcloud services enable bigquery.googleapis.com --project $PROJECT_ID

# Create a Pub/Sub topic for alerts
echo "Creating Pub/Sub topic for alerts..."
gcloud pubsub topics create gcp-guardrail-alerts --project $PROJECT_ID || true

echo "Setup complete!"
echo "Next steps:"
echo "1. Run the IAM analyzer: ./bin/iam-analyzer --project=$PROJECT_ID"
echo "2. Run the Terraform validator: ./bin/tf-validator --plan=<terraform-plan-file>"
echo "3. Start the log watcher: ./bin/log-watcher --project=$PROJECT_ID"
echo "4. Analyze service accounts: ./bin/sa-tracker --project=$PROJECT_ID"
echo "5. Start the Rego playground: ./bin/rego-playground" 