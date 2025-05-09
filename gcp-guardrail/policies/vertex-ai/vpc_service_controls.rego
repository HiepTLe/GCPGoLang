package vertex.vpc_service_controls

# This policy enforces that Vertex AI resources are protected by VPC Service Controls
# to prevent data exfiltration and enforce perimeter security

import input.planned_values as planned
import input.resource_changes as changes

# Default to deny
default deny = false

# Default empty message list
default messages = []

# List of Vertex AI resource types that should be inside service perimeters
vertex_resources = [
    "google_vertex_ai_dataset",
    "google_vertex_ai_model",
    "google_vertex_ai_endpoint",
    "google_vertex_ai_tensorboard",
    "google_vertex_ai_featurestore",
    "google_vertex_ai_metadata_store"
]

# Check that Vertex AI resources are within a VPC Service Control perimeter
violation[message] {
    # Filter for vertex AI resources
    resource := changes[_]
    resource.type == vertex_resources[_]
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check if the resource has network isolation configuration
    not after.network_isolation
    
    message := {
        "id": resource.address,
        "policy": "vertex.vpc_service_controls",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI resources must have network isolation enabled for VPC Service Controls compatibility",
        "remediation": "Set network_isolation = true and ensure the project is within a VPC Service Controls perimeter"
    }
}

# Check that Vertex AI resources are within a VPC Service Control perimeter
violation[message] {
    # Filter for Access Context Manager service perimeter resources
    resource := changes[_]
    resource.type == "google_access_context_manager_service_perimeter"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check if the AI Platform service is included in the perimeter
    after.restricted_services
    count(after.restricted_services) > 0
    
    # Check if AI Platform is in the list
    not contains_ai_platform(after.restricted_services)
    
    message := {
        "id": resource.address,
        "policy": "vertex.vpc_service_controls",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "VPC Service Control perimeters must include AI Platform service",
        "remediation": "Add 'aiplatform.googleapis.com' to the restricted_services list"
    }
}

# Helper function to check if AI Platform is in the restricted services list
contains_ai_platform(services) {
    service := services[_]
    service == "aiplatform.googleapis.com"
}

# Deny if any violations are found
deny {
    count(violation) > 0
    messages = violation
} 