# GCPGoLang Security Suite

A comprehensive cloud security governance toolkit for Google Cloud Platform using Go and Rego.

## Project Information
- **Project Name:** GCPGoLang
- **Project Number:** 652769711122
- **Project ID:** gcpgolang

## Latest Updates

- **Enhanced Terraform Validator**: Now provides detailed policy violation information with specific remediation guidance
- **Improved GitHub Actions Workflows**: Updated to use the latest Action versions and improved security reporting
- **Workload Identity Federation**: Robust authentication between GitHub Actions and GCP without service account keys
- **Severity-Based Validation**: Configurable validation thresholds to distinguish between warnings and errors
- **Dedicated Wiz Security Scanning**: Weekly comprehensive security scanning with customizable scan types and HTML reports

## Overview

GCPGoLang Security Suite is a collection of tools designed to enhance security governance in Google Cloud Platform environments. The toolkit helps identify security risks, enforce policies, monitor for threats, and ensure compliance. It integrates with and extends the gcp-guardrail project, following Go best practices with a modular architecture.

## Architecture

The project follows a modular architecture with:
- **Main CLI Application**: Centralized entry point with subcommands
- **Command Modules**: Separate packages for each functionality component
- **Core Services**: Shared functionality for GCP API interactions
- **Policy Engine**: Rego-based policy definitions and evaluation

## Components

### 1. IAM Analyzer
Analyzes GCP IAM policies to identify overly permissive permissions, policy violations, and security risks.

Features:
- Detects overprivileged accounts
- Identifies dangerous role combinations
- Flags external users with privileged roles
- Analyzes service account permissions
- Recommends least privilege adjustments

```bash
gcpgolang iam-analyzer --project=your-project-id
```

### 2. Service Account Tracker
Tracks and analyzes GCP Service Account usage patterns to identify unused accounts, over-permissioned accounts, and anomalous behavior.

Features:
- Identifies inactive service accounts
- Detects accounts with excess permissions
- Tracks key usage and rotation
- Monitors service account activity
- Generates usage reports in multiple formats

```bash
gcpgolang sa-tracker --project=your-project-id
```

### 3. Log Watcher
Monitors GCP audit logs for potential security threats and violations, detecting suspicious activities and generating alerts.

Features:
- Real-time log analysis
- Detection of suspicious access patterns
- Alerting for security violations
- Integration with Pub/Sub for notifications
- Customizable detection rules

```bash
gcpgolang log-watcher --project=your-project-id
```

### 4. Terraform Validator
Validates Terraform plans for GCP against security policies defined in Rego, checking for configuration issues, security risks, and policy violations before applying.

Features:
- Pre-deployment security validation
- Compliance checking against custom policies
- Detection of insecure configurations
- Integration with CI/CD pipelines
- Detailed violation reporting with remediation guidance
- Configurable severity thresholds for warnings vs. errors

```bash
# Basic validation
gcpgolang tf-validator --plan=path/to/terraform-plan.json

# With severity filtering (1=Low to 5=Blocker)
gcpgolang tf-validator --plan=path/to/terraform-plan.json --severity=2

# With custom failure threshold (only fail on high severity)
gcpgolang tf-validator --plan=path/to/terraform-plan.json --fail-threshold=4

# With custom policy directory
gcpgolang tf-validator --plan=path/to/terraform-plan.json --policy-dir=my-policies
```

Sample output:
```
Violation #1:
  Severity:      High
  Policy:        public_bucket_access
  Resource Type: google_storage_bucket
  Resource Name: my-non-compliant-bucket
  Issue:         Storage bucket does not have uniform bucket-level access enabled
  Remediation:   Set uniform_bucket_level_access = true in the bucket configuration
```

### 5. Misconfiguration Scanner
Scans GCP resources for security misconfigurations and integrates with Wiz for comprehensive vulnerability management.

Features:
- Detects GCP resource misconfigurations
- Categorizes findings by severity and resource type
- Provides actionable recommendations
- Integrates with Wiz for vulnerability data
- Generates comprehensive security reports

```bash
# Basic misconfiguration scan
gcpgolang misconfig-scanner --project=your-project-id

# With Wiz integration for vulnerability data
gcpgolang misconfig-scanner --project=your-project-id --wiz --wiz-client-id=YOUR_CLIENT_ID --wiz-client-secret=YOUR_CLIENT_SECRET

# Scan specific resource types
gcpgolang misconfig-scanner --project=your-project-id --scan-type=storage
```

## Policy Framework

The project uses Open Policy Agent (OPA) and Rego language for policy definition and enforcement:

- **IAM Policies**: Rules for proper access management
- **Network Policies**: Secure VPC and firewall configurations
- **Storage Policies**: Secure bucket and object configurations
- **Compute Policies**: VM and container security standards
- **Encryption Policies**: Data protection requirements

Policies can be customized to match organization-specific security requirements and compliance frameworks.

## GitHub Actions Integration

GCPGoLang integrates with GitHub Actions for continuous security monitoring of your GCP infrastructure:

### Workflow Features

- **Automated Scanning**: Runs on push, pull request, and scheduled intervals
- **Comprehensive Checks**: IAM analysis, service account auditing, and Terraform validation
- **Security Reports**: Generates detailed security reports as workflow artifacts
- **PR Integration**: Adds security findings as comments on pull requests
- **Audit History**: Maintains a record of security posture over time
- **Dedicated Wiz Scanning**: Weekly deep security scans with rich HTML reports and notification capabilities

### Enterprise-Grade Authentication with Workload Identity Federation

GCPGoLang uses Google Cloud's Workload Identity Federation for secure, key-less authentication from GitHub Actions. This approach follows cloud security best practices by:

1. **Eliminating service account keys**: No long-lived credentials to manage or risk exposing
2. **Using short-lived tokens**: Temporary OAuth tokens with limited scope and lifetime
3. **Fine-grained access control**: Precise permission boundaries for GitHub workflows
4. **Complete audit trail**: Every authentication and action is logged for security analysis
5. **Monitoring & alerting**: Suspicious activity detection with automated notifications

### Setup Using Terraform (Recommended)

We provide a complete Terraform configuration to set up Workload Identity Federation with proper RBAC and monitoring:

1. Navigate to the Terraform directory:
   ```bash
   cd terraform/github-actions
   ```

2. Copy and customize the variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your project details
   ```

3. Apply the Terraform configuration:
   ```bash
   terraform init
   terraform apply
   ```

4. Add the required secrets to your GitHub repository:
   - `WIF_PROVIDER`: Output from terraform (workload_identity_provider)
   - `SERVICE_ACCOUNT`: Output from terraform (service_account_email)
   - `GCP_PROJECT_ID`: Your Google Cloud project ID

For detailed instructions, see the [Workload Identity Federation README](terraform/github-actions/README.md).

### Manual Setup (Alternative)

If you prefer to set up Workload Identity Federation manually:

1. Create a Workload Identity Pool with a timestamp to avoid naming conflicts:
   ```bash
   TIMESTAMP=$(date +%m%d%H%M)
   gcloud iam workload-identity-pools create "github-pool-${TIMESTAMP}" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --display-name="GitHub Actions Pool"
   ```

2. Create a Workload Identity Provider:
   ```bash
   gcloud iam workload-identity-pools providers create-oidc "github-provider-${TIMESTAMP}" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --workload-identity-pool="github-pool-${TIMESTAMP}" \
     --display-name="GitHub Actions Provider" \
     --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
     --issuer-uri="https://token.actions.githubusercontent.com" \
     --attribute-condition="attribute.repository==\"YOUR_GITHUB_ORG/GCPGoLang\""
   ```

3. Create a service account and grant it permissions:
   ```bash
   gcloud iam service-accounts create "github-actions-sa" \
     --project="YOUR_PROJECT_ID"
   
   # Grant needed permissions (example)
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="serviceAccount:github-actions-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/iam.securityReviewer"
   ```

4. Allow the GitHub identity to impersonate the service account:
   ```bash
   gcloud iam service-accounts add-iam-policy-binding \
     "github-actions-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --project="YOUR_PROJECT_ID" \
     --role="roles/iam.workloadIdentityUser" \
     --member="principalSet://iam.googleapis.com/projects/YOUR_PROJECT_NUMBER/locations/global/workloadIdentityPools/github-pool-${TIMESTAMP}/attribute.repository/YOUR_GITHUB_ORG/GCPGoLang"
   ```

5. Add the required secrets to your GitHub repository as described above.

### Troubleshooting WIF Setup

If you encounter issues with Workload Identity Federation setup:

1. Run the information retrieval script:
   ```bash
   ./scripts/get-wif-info.sh
   ```

2. If resources are not appearing properly, run the cleanup script:
   ```bash
   ./scripts/cleanup-wif.sh
   ```

3. Create a new setup with a unique timestamp:
   ```bash
   ./scripts/manual-setup.sh
   ```

For detailed documentation, see [Workload Identity Federation Setup](docs/wif-setup.md).

## Complete Workflow

Here's a typical workflow for using GCPGoLang in your organization:

1. **Initial Setup**:
   - Install GCPGoLang and authenticate with GCP
   - Configure GitHub Actions (optional)
   - Customize Rego policies for your compliance requirements

2. **Continuous Security Monitoring**:
   - Schedule weekly IAM analysis scans
   - Monitor service account usage patterns
   - Watch audit logs for suspicious activity
   - Scan for misconfigurations and vulnerabilities

3. **Infrastructure Change Management**:
   - Validate Terraform plans before deployment
   - Block changes that violate security policies
   - Document exceptions and remediation plans

4. **Compliance Reporting**:
   - Generate periodic security reports
   - Track remediation progress
   - Document compliance status for audits

5. **Policy Refinement**:
   - Update policies based on new threats
   - Tune rules to reduce false positives
   - Add specific checks for your environment

## Integration with Wiz

GCPGoLang can integrate with the Wiz security platform to provide advanced vulnerability management:

### How It Works

1. The Misconfiguration Scanner detects GCP resource configuration issues
2. Wiz API integration fetches vulnerability data
3. Results are combined into a unified security report
4. Actionable recommendations are provided for both misconfigurations and vulnerabilities

### Dedicated Wiz Security Scanning

The project includes a dedicated GitHub Actions workflow for deeper Wiz security scanning:

- **Scheduled Weekly Scans**: Automatically runs every Monday at 2 AM
- **On-Demand Scanning**: Manually trigger scans via workflow dispatch
- **Customizable Scan Types**:
  - `full`: Complete scan of all resources and vulnerabilities
  - `critical-only`: Focus on critical and high severity issues
  - `compliance`: Validate against compliance frameworks (CIS, NIST, SOC2)
- **Rich Reporting**: Generates both JSON data and formatted HTML reports
- **Security Notifications**: Configurable notifications for critical findings

To run a manual scan:
1. Go to the Actions tab in your GitHub repository
2. Select "Wiz Security Scan" workflow
3. Click "Run workflow"
4. Choose your scan type and run

### Setup

1. Create a Wiz API client in your Wiz console:
   - Navigate to Settings > Service Accounts
   - Create a new service account with appropriate permissions
   - Generate client credentials

2. Add Wiz credentials to GitHub repository secrets:
   - `WIZ_CLIENT_ID`: Your Wiz API client ID
   - `WIZ_CLIENT_SECRET`: Your Wiz API client secret

3. Optionally configure notification channels:
   - `SECURITY_WEBHOOK_URL`: Webhook URL for security notifications (e.g., Slack, Teams, etc.)

## Getting Started

### Prerequisites
- Go 1.21 or later
- Google Cloud SDK (gcloud)
- Access to a GCP project

### Installation

1. Clone the repository:
```bash
git clone https://github.com/hieptle/GCPGoLang.git
cd GCPGoLang
```

2. Build the tools:
```bash
go build -o gcpgolang
```

3. Install the Google Cloud SDK if not already installed:
```bash
brew install --cask google-cloud-sdk  # On macOS with Homebrew
# OR
curl https://sdk.cloud.google.com | bash  # Other platforms
```

4. Initialize gcloud and authenticate:
```bash
gcloud init
gcloud auth application-default login
```

### Usage

Run the main application to see available commands:
```bash
./gcpgolang --help
```

## Project Setup with Terraform

GCPGoLang provides Terraform modules to automate the setup and configuration of your Google Cloud environment following security best practices.

### CI/CD for Infrastructure

All Terraform changes should be made through the GitHub Actions CI/CD pipeline:

1. **Make infrastructure changes** in the Terraform code
2. **Create a pull request** to propose the changes
3. **Review the Terraform plan** automatically added as a comment to your PR
4. **Merge the PR** to apply changes to the GCP environment

The workflow enforces:
- Consistent formatting and validation
- Pre-deployment planning and review
- Secure credential management
- Audit trail of all infrastructure changes

```bash
# To trigger an infrastructure deployment manually
# Go to Actions → Terraform Infrastructure Management → Run workflow
```

### GitHub Actions Authentication Setup

Before using the CI/CD pipeline, you need to set up authentication between GitHub Actions and GCP using Workload Identity Federation:

```bash
# Run the setup script with default values
./scripts/setup-github-auth.sh

# Or customize the setup
./scripts/setup-github-auth.sh --project-id=your-project-id --github-org=your-github-username
```

The script will output the values needed for GitHub repository secrets:
- `WIF_PROVIDER` - Workload Identity Provider resource name
- `SERVICE_ACCOUNT` - Service account email
- `GCP_PROJECT_ID` - Your GCP project ID

Add these secrets in your GitHub repository under Settings → Secrets and variables → Actions.

For detailed instructions, see [GitHub Secrets Setup](docs/github-secrets-setup.md).

### Project Infrastructure Setup

The `terraform/project-setup` module handles the core infrastructure setup:

1. **Enables Required APIs** - All necessary Google Cloud APIs are activated
2. **Configures Security Monitoring** - Sets up audit logs and alerting
3. **Establishes Security Foundations** - Creates resources needed for security tooling

```bash
cd terraform/project-setup
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your project details
terraform init
terraform apply
```

For detailed instructions, see the [Project Setup README](terraform/project-setup/README.md).

### GitHub Actions Authentication

The `terraform/github-actions` module configures secure authentication from GitHub Actions to Google Cloud using Workload Identity Federation:

```bash
cd terraform/github-actions
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your project details
terraform init
terraform apply
```

For details, see the [GitHub Actions Authentication README](terraform/github-actions/README.md).

## Development

### Project Structure
```
GCPGoLang/
├── main.go                  # Main entry point
├── go.mod                   # Go module definition
├── go.sum                   # Go dependency checksums
├── .github/                 # GitHub configuration
│   └── workflows/           # GitHub Actions workflows
├── gcp-guardrail/           # Core security framework
│   ├── cmd/                 # Command implementations
│   ├── pkg/                 # Core packages
│   │   ├── cmd/             # Command definitions
│   │   ├── gcp/             # GCP-specific logic
│   │   └── policy/          # Policy evaluation
│   └── policies/            # Rego policy definitions
├── docs/                    # Documentation
├── scripts/                 # Utility scripts
├── terraform/               # Terraform modules
│   ├── project-setup/       # Project foundation setup
│   └── github-actions/      # GitHub Actions authentication
└── examples/                # Usage examples
    ├── terraform/           # Terraform examples
    │   ├── plans/           # Example Terraform plans
    │   └── policies/        # Example Rego policies
    └── workflows/           # Example workflow files
```

### Adding New Components
1. Create a new package in `gcp-guardrail/pkg/cmd/`
2. Implement the command interface with `GetCommand()`
3. Add any supporting logic in `gcp-guardrail/pkg/gcp/`
4. Add relevant policies in `gcp-guardrail/policies/`
5. Register the command in `main.go`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
