# GCPGoLang Examples and GitHub Actions Integration

This directory contains examples demonstrating how to use GCPGoLang with GitHub Actions for continuous security scanning of your GCP infrastructure.

## GitHub Actions Integration

GCPGoLang can be integrated with GitHub Actions to automatically scan your GCP environment for security issues. This provides:

1. Continuous monitoring of your GCP security posture
2. Documentation of security issues as part of your CI/CD pipeline
3. Blocking unsafe infrastructure changes via policy checks
4. Audit history of security scans over time

### Required GitHub Secrets

To use the provided GitHub Actions workflow, you need to configure the following repository secrets:

- `GCP_PROJECT_ID`: Your Google Cloud project ID
- `WIF_PROVIDER`: The Workload Identity Federation provider URL
- `SERVICE_ACCOUNT`: The service account email address to use for authentication

### How to Setup GCP Workload Identity Federation

1. Create a Workload Identity Pool:
   ```bash
   gcloud iam workload-identity-pools create "github-actions-pool" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --display-name="GitHub Actions Pool"
   ```

2. Create a Workload Identity Provider:
   ```bash
   gcloud iam workload-identity-pools providers create-oidc "github-provider" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --workload-identity-pool="github-actions-pool" \
     --display-name="GitHub Provider" \
     --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
     --issuer-uri="https://token.actions.githubusercontent.com"
   ```

3. Create a Service Account:
   ```bash
   gcloud iam service-accounts create "github-actions-sa" \
     --project="YOUR_PROJECT_ID" \
     --display-name="GitHub Actions Service Account"
   ```

4. Grant the Service Account the necessary roles:
   ```bash
   gcloud projects add-iam-policy-binding "YOUR_PROJECT_ID" \
     --member="serviceAccount:github-actions-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/viewer"
     
   gcloud projects add-iam-policy-binding "YOUR_PROJECT_ID" \
     --member="serviceAccount:github-actions-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/iam.securityReviewer"
   ```

5. Allow the GitHub repository to impersonate the service account:
   ```bash
   gcloud iam service-accounts add-iam-policy-binding "github-actions-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --project="YOUR_PROJECT_ID" \
     --role="roles/iam.workloadIdentityUser" \
     --member="principalSet://iam.googleapis.com/projects/YOUR_PROJECT_NUMBER/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository/YOUR_GITHUB_USERNAME/YOUR_REPO_NAME"
   ```

6. Get the Workload Identity Provider resource name:
   ```bash
   gcloud iam workload-identity-pools providers describe "github-provider" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --workload-identity-pool="github-actions-pool" \
     --format="value(name)"
   ```

7. Add the provider name to your GitHub repository secrets as `WIF_PROVIDER` and the service account email as `SERVICE_ACCOUNT`.

## Workflow Outputs

The GitHub Actions workflow outputs:

1. **IAM Analysis Report**: Identifies IAM security issues like overly permissive roles
2. **Service Account Audit Report**: Lists unused or over-privileged service accounts
3. **Terraform Validation Reports**: Shows policy violations in Terraform plans
4. **Consolidated Security Report**: A summary of all security findings

## Customizing Rego Policies

The workflow uses the Rego policies in your repository. To customize the security checks:

1. Add or modify Rego policies in the `gcp-guardrail/policies/` directory
2. Create test cases in the `examples/terraform/plans/` directory
3. Test locally with `gcpgolang tf-validator --plan=your-plan.json`
4. Commit and push to trigger the GitHub Actions workflow

## Example Policy Explanation

The example policy in `examples/terraform/policies/storage_bucket.rego` demonstrates:

- Enforcing uniform bucket level access for security
- Requiring encryption on all storage buckets
- Recommending versioning for data protection

These policies will be evaluated against your Terraform plans and any violations will be reported in the GitHub Actions logs and as comments on pull requests. 