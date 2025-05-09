package vertex.responsible_ai

# This policy enforces responsible AI governance practices for machine learning models
# in Vertex AI, focusing on transparency, explainability, and governance.

import input.planned_values as planned
import input.resource_changes as changes

# Default to deny
default deny = false

# Default empty message list
default messages = []

# Check that ML models have proper metadata and documentation
violation[message] {
    # Filter for Vertex AI model resources
    resource := changes[_]
    resource.type == "google_vertex_ai_model"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check if metadata exists
    not after.metadata
    
    message := {
        "id": resource.address,
        "policy": "vertex.responsible_ai",
        "severity": "MEDIUM",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI models must have metadata for governance and documentation",
        "remediation": "Add metadata block with information about the model's purpose, training data, and evaluation metrics"
    }
}

# Check for explainability configuration
violation[message] {
    # Filter for Vertex AI model resources
    resource := changes[_]
    resource.type == "google_vertex_ai_model"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Missing explainability metadata
    after.metadata
    not after.explanation_spec
    
    message := {
        "id": resource.address,
        "policy": "vertex.responsible_ai",
        "severity": "LOW",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI models should have explainability configuration",
        "remediation": "Add explanation_spec block with feature attribution methods"
    }
}

# Check that Vertex AI pipelines include evaluation steps
violation[message] {
    # Filter for Vertex AI pipeline resources
    resource := changes[_]
    resource.type == "google_vertex_ai_pipeline_job"
    
    # Get "after" planned state
    after := resource.change.after
    
    # Check if pipeline_spec exists and contains evaluation components
    not contains(after.pipeline_spec, "evaluation")
    not contains(after.pipeline_spec, "metrics")
    
    message := {
        "id": resource.address,
        "policy": "vertex.responsible_ai",
        "severity": "MEDIUM",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Vertex AI pipelines should include model evaluation components",
        "remediation": "Include evaluation steps in your pipeline definition that measure model performance and fairness metrics"
    }
}

# Check for model card creation 
violation[message] {
    # Filter for Vertex AI model deployment
    resource := changes[_]
    resource.type == "google_vertex_ai_endpoint"
    
    # Get "after" planned state
    after := resource.change.after
    
    # If deploying a model to an endpoint, model card should exist
    after.deployed_models
    count(after.deployed_models) > 0
    
    # Check model card metadata exists
    deployed_model := after.deployed_models[_]
    not contains(deployed_model.model, "model_card")
    
    message := {
        "id": resource.address,
        "policy": "vertex.responsible_ai",
        "severity": "LOW",
        "resource_type": resource.type,
        "resource_name": resource.name,
        "message": "Models deployed to endpoints should have model cards for governance",
        "remediation": "Create a model card with information about intended use, limitations, and ethical considerations before deployment"
    }
}

# Deny if any violations are found
deny {
    count(violation) > 0
    messages = violation
} 