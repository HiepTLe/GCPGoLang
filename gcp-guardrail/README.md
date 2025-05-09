# GCP GuardRail - Cloud Security Governance Toolkit

A modular Golang-based cloud governance toolkit for GCP, combining policy validation, security automation, IAM audit analysis, and OPA-based policy enforcement. Designed for enterprise use at scale.

## Features

- **IAM Policy Analyzer**: Analysis tool for GCP IAM policies to enforce least privilege
- **Terraform Compliance Validator**: Static analysis of Terraform plans against security policies
- **GCP Log Sink Threat Detector**: Real-time monitoring of GCP audit logs for security threats
- **Service Account Usage Tracker**: Monitor and analyze service account usage patterns
- **Rego Playground**: Web-based interface for testing and developing Rego policies
- **OPA Admission Controller**: Kubernetes admission controller for enforcing policies
- **Misconfiguration Scanner**: Detect and report security misconfigurations in GCP resources

## Architecture

```
gcp-guardrail/
├── cmd/                          # Command-line tools
│   ├── iam-analyzer/             # IAM policy analysis CLI
│   ├── tf-validator/             # Terraform compliance scanner
│   ├── log-watcher/              # Audit log threat detector
│   └── sa-tracker/               # Service account usage analyzer
├── pkg/                          # Shared packages
│   ├── gcp/                      # GCP API helpers
│   ├── rego/                     # Rego policy wrappers
│   ├── terraform/                # Terraform plan parser
│   └── alerts/                   # Notification/email/Slack integrations
├── policies/                     # Rego policies
│   ├── iam/                      # IAM policies
│   ├── terraform/                # Terraform policies
│   └── logging/                  # Logging/audit policies
├── playground/                   # Web UI for testing Rego policies
├── controller/                   # OPA-based K8s admission controller
├── examples/                     # Sample inputs, use cases
└── scripts/                      # Setup scripts
```

## Getting Started

### Prerequisites

- Go 1.21+
- Google Cloud SDK
- Terraform (for IaC validation features)
- Docker (for running the Rego playground)
- Kubernetes cluster (for admission controller)

### Installation

1. Clone the repository:
   ```
   git clone github.com/hieptle/gcp-guardrail
   cd gcp-guardrail
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build the tools:
   ```
   go build -o bin/ ./cmd/...
   ```

### Usage Examples

#### IAM Analyzer

```bash
./bin/iam-analyzer --project=my-gcp-project --report-format=json
```

#### Terraform Validator

```bash
./bin/tf-validator --plan=terraform-plan.json --policies=policies/terraform
```

#### Log Watcher

```bash
./bin/log-watcher --project=my-gcp-project --sink=my-log-sink
```

#### Service Account Tracker

```bash
./bin/sa-tracker --project=my-gcp-project --days=30
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 