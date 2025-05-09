package vertex.access_restrictions

# This policy enforces proper IAM access restrictions for Vertex AI resources
# to prevent overly permissive access to sensitive ML models and data

import input.planned_values as planned
import input.resource_changes as changes

# Default to deny
default deny = false

# Default empty message list
default messages = []

# List of Vertex AI resource IAM bindings that we want to check
vertex_iam_resources = [
    "google_vertex_ai_dataset_iam_binding",
    "google_vertex_ai_dataset_iam_member",
    "google_vertex_ai_model_iam_binding",
    "google_vertex_ai_model_iam_member",
    "google_vertex_ai_endpoint_iam_binding",
    "google_vertex_ai_endpoint_iam_member",
    "google_vertex_ai_featurestore_iam_binding",
    "google_vertex_ai_featurestore_iam_member"
]

# High-privileged roles that should be restricted
sensitive_roles = [
    "roles/aiplatform.admin",
    "roles/aiplatform.user",
    "roles/aiplatform.modelUser"
]

# Check for overly permissive IAM bindings (allUsers, allAuthenticatedUsers)
violation[message] {
    # Filter for Vertex AI IAM resources
    resource := changes[_]
    resource.type == vertex_iam_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check for public access in IAM bindings
    resource.type == concat("", [_, "_iam_binding"])
    public_member := after.members[_]
    public_member == "allUsers" or public_member == "allAuthenticatedUsers"
    
    message := {
        "id": resource.address,
        "policy": "vertex.access_restrictions",
        "severity": "CRITICAL",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI resources cannot have public IAM access (allUsers or allAuthenticatedUsers)",
        "remediation": "Remove public members from IAM bindings"
    }
}

# Check for overly permissive IAM member assignments
violation[message] {
    # Filter for Vertex AI IAM resources
    resource := changes[_]
    resource.type == vertex_iam_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check for public access in IAM member
    resource.type == concat("", [_, "_iam_member"])
    after.member == "allUsers" or after.member == "allAuthenticatedUsers"
    
    message := {
        "id": resource.address,
        "policy": "vertex.access_restrictions",
        "severity": "CRITICAL",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI resources cannot have public IAM access (allUsers or allAuthenticatedUsers)",
        "remediation": "Use specific user, group, or service account identities instead of public members"
    }
}

# Check for sensitive roles assigned to service accounts outside our organization
violation[message] {
    # Filter for Vertex AI IAM resources
    resource := changes[_]
    resource.type == vertex_iam_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check if a sensitive role is being assigned
    after.role == sensitive_roles[_]
    
    # For member assignments
    resource.type == concat("", [_, "_iam_member"])
    startswith(after.member, "serviceAccount:")
    not contains(after.member, ".gserviceaccount.com")
    
    message := {
        "id": resource.address,
        "policy": "vertex.access_restrictions",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Sensitive Vertex AI roles should only be assigned to service accounts within the organization",
        "remediation": "Restrict access to service accounts within your organization's domain"
    }
}

# Check for sensitive roles in bindings
violation[message] {
    # Filter for Vertex AI IAM resources
    resource := changes[_]
    resource.type == vertex_iam_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # For binding assignments with sensitive roles
    resource.type == concat("", [_, "_iam_binding"])
    after.role == sensitive_roles[_]
    
    # Check each member in binding
    member := after.members[_]
    startswith(member, "serviceAccount:")
    not contains(member, ".gserviceaccount.com")
    
    message := {
        "id": resource.address,
        "policy": "vertex.access_restrictions",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Sensitive Vertex AI roles should only be assigned to service accounts within the organization",
        "remediation": "Restrict access to service accounts within your organization's domain"
    }
}

# Deny if any violations are found
deny {
    count(violation) > 0
    messages = violation
} 