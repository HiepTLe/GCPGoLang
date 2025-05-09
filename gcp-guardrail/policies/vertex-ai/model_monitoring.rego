package vertex.model_monitoring

# This policy enforces that deployed machine learning models have
# monitoring enabled to detect model drift and data quality issues

import input.planned_values as planned
import input.resource_changes as changes

# Default to deny
default deny = false

# Default empty message list
default messages = []

# Check all model deployment resources
violation[message] {
    # Filter for model deployment resources
    resource := changes[_]
    resource.type == "google_vertex_ai_model_deployment_monitoring_job"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Define a list of required monitoring tasks
    required_tasks = ["feature_drift", "prediction_drift"]
    
    # Check if monitoring_config exists
    not after.model_deployment_monitoring_job_config
    
    message := {
        "id": resource.address,
        "policy": "vertex.model_monitoring",
        "severity": "MEDIUM",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI model deployments must have monitoring configured",
        "remediation": "Add model_deployment_monitoring_job_config block with monitoring settings"
    }
}

# Check that monitoring includes drift detection
violation[message] {
    # Filter for model deployment resources
    resource := changes[_]
    resource.type == "google_vertex_ai_model_deployment_monitoring_job"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Has config but missing drift detection settings
    after.model_deployment_monitoring_job_config
    not after.model_deployment_monitoring_job_config.model_monitoring_objective.prediction_drift_detection_config
    not after.model_deployment_monitoring_job_config.model_monitoring_objective.feature_drift_detection_config
    
    message := {
        "id": resource.address,
        "policy": "vertex.model_monitoring",
        "severity": "MEDIUM",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI model monitoring must include either prediction or feature drift detection",
        "remediation": "Configure model_monitoring_objective with drift detection settings"
    }
}

# Check model monitoring alert configuration
violation[message] {
    # Filter for model deployment resources
    resource := changes[_]
    resource.type == "google_vertex_ai_model_deployment_monitoring_job"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Has monitoring config but missing alerting
    after.model_deployment_monitoring_job_config
    not after.model_deployment_monitoring_job_config.notification_spec
    
    message := {
        "id": resource.address,
        "policy": "vertex.model_monitoring",
        "severity": "LOW",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI model monitoring should have alerting configured",
        "remediation": "Add notification_spec block to enable alerts for drift detection"
    }
}

# Deny if any violations are found
deny {
    count(violation) > 0
    messages = violation
} 