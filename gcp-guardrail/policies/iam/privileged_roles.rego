package gcp.iam

# List of high-risk roles that should be carefully monitored
# and assigned only when necessary
high_risk_roles = [
    "roles/owner",
    "roles/editor",
    "roles/iam.securityAdmin", 
    "roles/iam.serviceAccountAdmin",
    "roles/compute.admin",
    "roles/storage.admin",
    "roles/bigquery.admin",
    "roles/resourcemanager.projectIamAdmin",
    "roles/logging.admin",
    "roles/pubsub.admin"
]

# Deny high-risk roles assigned directly to users
# Allow for service accounts with proper naming convention
deny[msg] {
    binding := input.bindings[_]
    role := binding.role
    high_risk_roles[_] == role
    
    member := binding.members[_]
    startswith(member, "user:")
    
    msg := sprintf("High-risk role '%s' assigned to user '%s'", [role, member])
}

# Detect service accounts with owner/editor roles
deny[msg] {
    binding := input.bindings[_]
    role := binding.role
    role == "roles/owner" or role == "roles/editor"
    
    member := binding.members[_]
    startswith(member, "serviceAccount:")
    
    msg := sprintf("Critical role '%s' assigned to service account '%s'", [role, member])
}

# Detect external users with privileged roles
deny[msg] {
    binding := input.bindings[_]
    role := binding.role
    high_risk_roles[_] == role
    
    member := binding.members[_]
    startswith(member, "user:")
    not endswith(member, "@mycompany.com")
    
    msg := sprintf("High-risk role '%s' assigned to external user '%s'", [role, member])
}

# Count number of principals with owner role
owner_principals[principal] {
    binding := input.bindings[_]
    binding.role == "roles/owner"
    principal := binding.members[_]
}

# Warning if too many owners
warn[msg] {
    owners := count(owner_principals)
    owners > 3
    
    msg := sprintf("Too many principals (%d) with Owner role", [owners])
} 