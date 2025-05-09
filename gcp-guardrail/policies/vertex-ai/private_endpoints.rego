package vertex.private_endpoints

# This policy enforces that Vertex AI endpoints use private connectivity
# rather than public endpoints for enhanced security

import input.planned_values as planned
import input.resource_changes as changes

# Default to deny
default deny = false

# Default empty message list
default messages = []

# Check all Vertex AI endpoint resources
violation[message] {
    # Filter for vertex AI endpoint resources
    resource := changes[_]
    resource.type == "google_vertex_ai_endpoint"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check if network config exists and enforces private service access
    not after.network
    
    # If missing network config or not using private service networking
    message := {
        "id": resource.address,
        "policy": "vertex.private_endpoints",
        "severity": "HIGH",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI endpoints must use private connectivity via VPC network peering",
        "remediation": "Configure the endpoint with network = \"projects/{project_id}/global/networks/{network}\" to use private connectivity"
    }
}

# Allow only endpoints with proper private connectivity
deny {
    count(violation) > 0
    messages = violation
} 