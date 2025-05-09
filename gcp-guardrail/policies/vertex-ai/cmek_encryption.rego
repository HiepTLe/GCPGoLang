package vertex.cmek_encryption

# This policy enforces that Vertex AI resources use Customer-Managed Encryption Keys (CMEK)
# instead of Google-managed keys for enhanced security and compliance

import input.planned_values as planned
import input.resource_changes as changes

# Default to deny
default deny = false

# Default empty message list
default messages = []

# List of Vertex AI resource types that should use CMEK
vertex_resources = [
    "google_vertex_ai_dataset",
    "google_vertex_ai_model",
    "google_vertex_ai_tensorboard",
    "google_vertex_ai_featurestore",
    "google_vertex_ai_metadata_store"
]

# Check all Vertex AI resources for CMEK encryption
violation[message] {
    # Filter for vertex AI resources
    resource := changes[_]
    resource.type == vertex_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # Missing encryption_spec or not using a KMS key
    not after.encryption_spec
    
    message := {
        "id": resource.address,
        "policy": "vertex.cmek_encryption",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI resources must use Customer-Managed Encryption Keys (CMEK)",
        "remediation": "Add encryption_spec block with a KMS key reference"
    }
}

# Also check for empty encryption_spec or missing KMS key
violation[message] {
    # Filter for vertex AI resources
    resource := changes[_]
    resource.type == vertex_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # Has encryption_spec but missing kms_key_name
    after.encryption_spec
    not after.encryption_spec.kms_key_name
    
    message := {
        "id": resource.address,
        "policy": "vertex.cmek_encryption",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI resources with encryption_spec must specify a kms_key_name",
        "remediation": "Specify a valid KMS key in the encryption_spec.kms_key_name field"
    }
}

# Deny if any violations are found
deny {
    count(violation) > 0
    messages = violation
} 