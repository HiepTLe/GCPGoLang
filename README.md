# GCPGoLang Security Suite

A comprehensive cloud security governance toolkit for Google Cloud Platform using Go and Rego.

## Project Information
- **Project Name:** GCPGoLang
- **Project Number:** 652769711122
- **Project ID:** gcpgolang

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
- Detailed violation reporting

```bash
gcpgolang tf-validator --plan=path/to/terraform-plan.json
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

### Setup

1. Copy the workflow file to your repository:
   ```bash
   mkdir -p .github/workflows
   cp examples/workflows/gcpgolang.yml .github/workflows/
   ```

2. Set up GCP authentication using the provided script:
   ```bash
   ./scripts/setup-github-actions.sh YOUR_PROJECT_ID
   ```

3. Add the required secrets to your GitHub repository:
   - `GCP_PROJECT_ID`: Your Google Cloud project ID
   - `WIF_PROVIDER`: The Workload Identity Federation provider URL
   - `SERVICE_ACCOUNT`: The service account email address

For detailed instructions, see [GitHub Actions Integration](docs/github-actions-integration.md).

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

1. The Misconfiguration Scanner detects GCP configuration issues
2. Wiz API integration fetches vulnerability data
3. Results are combined into a unified security report
4. Actionable recommendations are provided for both misconfigurations and vulnerabilities

### Setup

1. Create a Wiz API client in your Wiz console:
   - Navigate to Settings > Service Accounts
   - Create a new service account with appropriate permissions
   - Generate client credentials

2. Run the scanner with Wiz integration:
   ```bash
   gcpgolang misconfig-scanner --project=your-project-id --wiz \
     --wiz-client-id=YOUR_CLIENT_ID \
     --wiz-client-secret=YOUR_CLIENT_SECRET
   ```

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
