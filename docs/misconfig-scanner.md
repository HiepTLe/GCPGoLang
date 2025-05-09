# Misconfiguration Scanner with Wiz Integration

The Misconfiguration Scanner is a powerful tool in the GCPGoLang suite that detects security misconfigurations in your GCP environment and can integrate with Wiz for comprehensive vulnerability management.

## Overview

This scanner analyzes your GCP resources to identify security issues such as:
- Public storage buckets
- Overly permissive firewall rules
- Unencrypted data
- Excessive IAM permissions
- And more

When integrated with Wiz, it also provides vulnerability information from the Wiz platform, creating a unified security view.

## Basic Usage

```bash
# Basic scan of all resource types
gcpgolang misconfig-scanner --project=your-project-id

# Output in JSON format
gcpgolang misconfig-scanner --project=your-project-id --report-format=json

# Save results to a file
gcpgolang misconfig-scanner --project=your-project-id --output=results.json

# Scan only specific resource types
gcpgolang misconfig-scanner --project=your-project-id --scan-type=storage
```

## Wiz Integration

The scanner can integrate with the Wiz security platform to include vulnerability information alongside misconfiguration data.

### Prerequisites

1. A Wiz account with API access
2. API credentials (Client ID and Client Secret)

### Setup Wiz Integration

1. Create API credentials in Wiz:
   - Go to your Wiz Admin Console
   - Navigate to Settings > Service Accounts
   - Create a new service account
   - Grant appropriate permissions (read-only access is sufficient)
   - Generate and save credentials

2. Run the scanner with Wiz integration:

```bash
gcpgolang misconfig-scanner --project=your-project-id --wiz \
  --wiz-client-id=YOUR_CLIENT_ID \
  --wiz-client-secret=YOUR_CLIENT_SECRET
```

## Command Options

| Option | Description | Default |
|--------|-------------|---------|
| `--project` | GCP project ID | Required |
| `--scan-type` | Resource types to scan (all, storage, compute, network, iam) | all |
| `--report-format` | Output format (text, json, csv) | text |
| `--output` | Output file path | (stdout) |
| `--wiz` | Enable Wiz integration | false |
| `--wiz-client-id` | Wiz API Client ID | |
| `--wiz-client-secret` | Wiz API Client Secret | |
| `--verbose` | Enable verbose output | false |

## Output Format

### Sample Text Output

```
GCP Misconfiguration Scan Results
Project: your-project-id
Scan Time: 2023-09-15T12:34:56Z

Total Issues Found: 6
  CRITICAL: 1
  HIGH: 2
  MEDIUM: 2
  LOW: 1

GCP Misconfigurations:
1. [CRITICAL] default-allow-all: Overly permissive firewall rule (0.0.0.0/0)
   Resource: compute.googleapis.com/Firewall
   Recommendation: Restrict firewall rules to specific IP ranges

2. [HIGH] example-bucket: Public access enabled
   Resource: storage.googleapis.com/Bucket
   Recommendation: Configure uniform bucket-level access and remove public access

3. [HIGH] service-account-1: Service account has owner role
   Resource: iam.googleapis.com/ServiceAccount
   Recommendation: Follow principle of least privilege and assign more specific roles

Wiz Vulnerabilities:
1. [CRITICAL] frontend-app: CVE-2023-1234
   Description: Critical vulnerability in container image
   First Seen: 2023-09-13T12:34:56Z
   Remediation: Update to latest version
   CVE: CVE-2023-1234

2. [MEDIUM] frontend-lb: Outdated TLS Configuration
   Description: Load balancer using outdated TLS configuration
   First Seen: 2023-09-12T12:34:56Z
   Remediation: Update TLS configuration to use TLS 1.2+
```

### JSON Output Structure

```json
{
  "project_id": "your-project-id",
  "scan_time": "2023-09-15T12:34:56Z",
  "misconfigurations": [
    {
      "resource_type": "storage.googleapis.com/Bucket",
      "resource_name": "example-bucket",
      "resource_id": "projects/your-project-id/buckets/example-bucket",
      "issue": "Public access enabled",
      "severity": "HIGH",
      "recommendation": "Configure uniform bucket-level access and remove public access",
      "timestamp": "2023-09-15T12:34:56Z",
      "category": "Storage"
    }
  ],
  "wiz_vulnerabilities": [
    {
      "id": "wiz-vuln-1",
      "name": "CVE-2023-1234",
      "description": "Critical vulnerability in container image",
      "severity": "CRITICAL",
      "resourceName": "frontend-app",
      "resourceType": "Container",
      "firstSeen": "2023-09-13T12:34:56Z",
      "status": "OPEN",
      "remediation": "Update to latest version",
      "cve": "CVE-2023-1234"
    }
  ],
  "total_issues": 2,
  "severity_counts": {
    "CRITICAL": 1,
    "HIGH": 1,
    "MEDIUM": 0,
    "LOW": 0
  }
}
```

## GitHub Actions Integration

The GCPGoLang GitHub Actions workflow automatically includes the misconfiguration scanner. To enable Wiz integration in CI/CD, add the following secrets to your GitHub repository:

- `WIZ_CLIENT_ID`: Your Wiz API Client ID
- `WIZ_CLIENT_SECRET`: Your Wiz API Client Secret

## Examples

### Scan Compute Resources

```bash
gcpgolang misconfig-scanner --project=your-project-id --scan-type=compute
```

### Generate JSON Report with Wiz Integration

```bash
gcpgolang misconfig-scanner --project=your-project-id --wiz \
  --wiz-client-id=YOUR_CLIENT_ID \
  --wiz-client-secret=YOUR_CLIENT_SECRET \
  --report-format=json --output=security-report.json
```

### Scan Multiple Projects

```bash
#!/bin/bash
for project in project1 project2 project3; do
  echo "Scanning $project..."
  gcpgolang misconfig-scanner --project=$project --output=$project-report.json
done
``` 