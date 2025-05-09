# Integrating GCPGoLang with GitHub Actions

This guide explains how to integrate GCPGoLang with GitHub Actions for continuous security monitoring of your Google Cloud Platform environment.

## Benefits of GitHub Actions Integration

1. **Continuous Security Monitoring**: Automatically scan your GCP environment on every change
2. **Security as Code**: Manage security policies alongside your infrastructure code
3. **Pipeline Integration**: Block insecure changes before they reach production
4. **Documentation**: Maintain a history of security posture over time
5. **Visibility**: Make security findings visible to all developers, not just security teams

## Setup Process

### 1. Create the GitHub Actions Workflow

Create the file `.github/workflows/gcpgolang.yml` in your repository with the content provided in the `examples` directory. This workflow will:

- Build the GCPGoLang tools
- Run IAM policy analysis
- Audit service account usage
- Validate Terraform plans against security policies
- Generate a consolidated security report

### 2. Set Up GCP Authentication

GCPGoLang uses GCP Workload Identity Federation for secure authentication from GitHub Actions. To set this up:

1. Run the setup script:
   ```bash
   ./scripts/setup-github-actions.sh YOUR_GCP_PROJECT_ID
   ```

2. Add the provided secrets to your GitHub repository:
   - `GCP_PROJECT_ID`
   - `WIF_PROVIDER`
   - `SERVICE_ACCOUNT`

Alternatively, follow the manual setup instructions in the `examples/README.md` file.

## Rego Policy Requirements

GCPGoLang uses Open Policy Agent (OPA) Rego policies for security evaluation. To work properly with GitHub Actions, your policies should:

### 1. Follow the Package Structure

Policies should be organized in a consistent package structure:

```
policies/
├── iam/              # IAM-related policies
├── compute/          # Compute Engine policies
├── storage/          # Cloud Storage policies
├── network/          # VPC and networking policies
├── kubernetes/       # GKE policies
└── terraform/        # Terraform-specific policies
```

### 2. Use the Expected Rule Format

Policies should define `deny` and `warn` rules:

```rego
package gcp.storage

# Rule to deny non-compliant resources
deny[msg] {
    # Logic to detect violations
    msg := sprintf("Violation: %v", [details])
}

# Rule to warn about potential issues
warn[msg] {
    # Logic to detect warnings
    msg := sprintf("Warning: %v", [details])
}
```

### 3. Provide Clear Violation Messages

Each violation should include:
- What resource has the issue
- What the issue is
- How to fix it

Example:
```rego
msg := sprintf("Storage bucket '%s' does not have uniform bucket level access enabled. Enable it to improve security.", [bucket.name])
```

### 4. Handle Input Formats Correctly

For IAM analysis, input will be in the format:
```json
{
  "bindings": [
    {
      "role": "roles/owner",
      "members": ["user:example@example.com"]
    }
  ]
}
```

For Terraform validation, input will be the Terraform plan in JSON format:
```json
{
  "resource_changes": [
    {
      "type": "google_storage_bucket",
      "change": {
        "after": {
          "name": "example-bucket",
          "uniform_bucket_level_access": false
        }
      }
    }
  ]
}
```

## GitHub Actions Workflow Output

The GCPGoLang GitHub Actions workflow generates several outputs:

### 1. Console Logs

During workflow execution, policy evaluation results are printed to the console logs:

```
Analyzing IAM policies for project my-project
Found 3 policy violations:
HIGH: Owner role assigned to user@example.com
MEDIUM: External user user@external.com has Editor access
LOW: Too many principals (5) with Owner role
```

### 2. Artifact Reports

The workflow generates and uploads artifact reports:
- `iam-analysis-report.json`: IAM policy analysis results
- `sa-report.json`: Service account audit results
- `*-report.json`: Terraform validation results for each plan
- `security-report.md`: Consolidated Markdown report of all findings

### 3. Pull Request Comments

For pull requests, the workflow automatically comments with a summary of all security findings:

```markdown
# GCPGoLang Security Scan Results
## Scan Date: Mon Sep 4 12:34:56 UTC 2023

## IAM Analysis
- HIGH: Owner role assigned to user@example.com
- MEDIUM: External user user@external.com has Editor access

## Service Account Analysis
- Unused account: unused-sa@my-project.iam.gserviceaccount.com
- Over-privileged: admin-sa@my-project.iam.gserviceaccount.com

## Terraform Validation
### Plan: storage_example
- ERROR: Storage bucket 'my-non-compliant-bucket' should have uniform bucket level access enabled for security
- ERROR: Storage bucket 'my-non-compliant-bucket' should have encryption configured
- WARNING: Storage bucket 'my-non-compliant-bucket' should have versioning enabled for data protection
```

## Customization Options

### Schedule Frequency

You can adjust the schedule in the GitHub Actions workflow:

```yaml
schedule:
  - cron: '0 0 * * 0'  # Weekly on Sundays
  # or
  - cron: '0 0 * * *'  # Daily at midnight
```

### Adding Custom Tests

To add custom security tests:

1. Create new Rego policies in the `policies/` directory
2. Add test cases in the `examples/terraform/plans/` directory
3. Update the workflow file if necessary for new test types

### Configuring Notifications

To add Slack notifications for security findings:

```yaml
- name: Notify Slack
  if: always()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    fields: repo,message,commit,author,action,eventName,workflow
    text: 'GCPGoLang Security Scan Results'
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:
- Verify the service account has the necessary permissions
- Check that the Workload Identity Federation is set up correctly
- Ensure all GitHub secrets are properly set

### Policy Evaluation Failures

If policies aren't evaluating as expected:
- Test the policies locally with the Rego playground
- Check input data formats match what the policies expect
- Review logs for parsing or syntax errors

### GitHub Actions Workflow Failures

For workflow failures:
- Check the job logs for specific error messages
- Verify dependencies and versions are correct
- Make sure any required files or directories exist 