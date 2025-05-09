package terraform.storage

# Deny public access to storage buckets
deny[msg] {
    resource := input.resource_changes[_]
    resource.type == "google_storage_bucket"
    resource.change.after.uniform_bucket_level_access == false
    
    msg := sprintf("Storage bucket '%s' should have uniform bucket level access enabled for security", [resource.change.after.name])
}

# Enforce encryption on storage buckets
deny[msg] {
    resource := input.resource_changes[_]
    resource.type == "google_storage_bucket"
    not resource.change.after.encryption
    
    msg := sprintf("Storage bucket '%s' should have encryption configured", [resource.change.after.name])
}

# Recommend versioning for important buckets
warn[msg] {
    resource := input.resource_changes[_]
    resource.type == "google_storage_bucket"
    not resource.change.after.versioning
    
    msg := sprintf("Storage bucket '%s' should have versioning enabled for data protection", [resource.change.after.name])
} 