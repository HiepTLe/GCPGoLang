package gcp.terraform.storage

import data.gcp.terraform.util as util

# Deny public buckets - check uniform bucket level access is enabled
deny[msg] {
    bucket := input.resources.google_storage_bucket[_]
    not bucket.attributes.uniform_bucket_level_access
    
    msg := sprintf("Storage bucket '%s' does not have uniform bucket level access enabled", [bucket.name])
}

# Deny buckets without versioning
deny[msg] {
    bucket := input.resources.google_storage_bucket[_]
    not bucket.attributes.versioning[0].enabled
    
    msg := sprintf("Storage bucket '%s' does not have versioning enabled", [bucket.name])
}

# Deny buckets with default logging disabled
deny[msg] {
    bucket := input.resources.google_storage_bucket[_]
    not bucket.attributes.logging
    
    msg := sprintf("Storage bucket '%s' does not have logging enabled", [bucket.name])
}

# Deny buckets without encryption
deny[msg] {
    bucket := input.resources.google_storage_bucket[_]
    not bucket.attributes.encryption
    
    msg := sprintf("Storage bucket '%s' does not have CMEK encryption configured", [bucket.name])
}

# Check for public access prevention
deny[msg] {
    bucket := input.resources.google_storage_bucket[_]
    not bucket.attributes.public_access_prevention
    
    msg := sprintf("Storage bucket '%s' does not have public access prevention set to 'enforced'", [bucket.name])
}

# Check for lifecycle rules (recommended)
warn[msg] {
    bucket := input.resources.google_storage_bucket[_]
    not bucket.attributes.lifecycle_rule
    
    msg := sprintf("Storage bucket '%s' does not have lifecycle rules configured", [bucket.name])
}

# Check for proper naming convention
warn[msg] {
    bucket := input.resources.google_storage_bucket[_]
    name := bucket.attributes.name
    not startswith(name, "gcp-")
    
    msg := sprintf("Storage bucket '%s' does not follow naming convention (should start with 'gcp-')", [name])
}

# Calculate compliance score
compliance_score = util.calc_compliance_score {
    true
}

# Helper to calculate compliance score
util = {
    "calc_compliance_score": score
} {
    total_denies := count(deny)
    total_warns := count(warn)
    total_checks := 7
    
    score := (total_checks - total_denies - (total_warns * 0.5)) / total_checks * 100
} 