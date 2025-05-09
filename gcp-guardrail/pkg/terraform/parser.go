package terraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Resource represents a Terraform resource
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Change represents a resource change in a Terraform plan
type Change struct {
	Resource    Resource               `json:"resource"`
	Action      string                 `json:"action"` // create, update, delete
	Before      map[string]interface{} `json:"before"`
	After       map[string]interface{} `json:"after"`
	Replacements []string               `json:"replacements,omitempty"`
}

// Plan represents a parsed Terraform plan
type Plan struct {
	FormatVersion    string    `json:"format_version"`
	TerraformVersion string    `json:"terraform_version"`
	Variables        map[string]interface{} `json:"variables"`
	ResourceChanges  []Change  `json:"resource_changes"`
	OutputChanges    map[string]interface{} `json:"output_changes"`
}

// Parser handles Terraform plan parsing
type Parser struct {
	planPath string
}

// NewParser creates a new Terraform plan parser
func NewParser(planPath string) *Parser {
	return &Parser{
		planPath: planPath,
	}
}

// Parse parses a Terraform plan file (JSON format)
func (p *Parser) Parse() (*Plan, error) {
	// Read the plan file
	data, err := ioutil.ReadFile(p.planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	// Parse the JSON
	var plan Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w", err)
	}

	return &plan, nil
}

// GetGCPResources filters the plan for GCP resources only
func (p *Plan) GetGCPResources() []Resource {
	var resources []Resource

	for _, change := range p.ResourceChanges {
		// Filter for GCP resources only
		if isGCPResource(change.Resource.Type) {
			resources = append(resources, change.Resource)
		}
	}

	return resources
}

// GetResourcesByAction filters resources by the specified action (create, update, delete)
func (p *Plan) GetResourcesByAction(action string) []Change {
	var changes []Change

	for _, change := range p.ResourceChanges {
		if change.Action == action {
			changes = append(changes, change)
		}
	}

	return changes
}

// isGCPResource determines if a resource type belongs to GCP
func isGCPResource(resourceType string) bool {
	gcpPrefixes := []string{
		"google_",
		"google-beta_",
	}

	for _, prefix := range gcpPrefixes {
		if len(resourceType) > len(prefix) && resourceType[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

// ConvertPlanToOPAInput converts the Terraform plan to a format suitable for OPA evaluation
func (p *Plan) ConvertPlanToOPAInput() map[string]interface{} {
	input := make(map[string]interface{})
	input["terraform_version"] = p.TerraformVersion
	
	// Extract GCP resources by type
	resourcesByType := make(map[string][]map[string]interface{})
	
	for _, change := range p.ResourceChanges {
		if isGCPResource(change.Resource.Type) && change.Action != "delete" {
			// Use the "after" state for resource evaluation
			resourceData := map[string]interface{}{
				"name":       change.Resource.Name,
				"attributes": change.After,
			}
			
			resourceType := change.Resource.Type
			if _, exists := resourcesByType[resourceType]; !exists {
				resourcesByType[resourceType] = []map[string]interface{}{}
			}
			
			resourcesByType[resourceType] = append(resourcesByType[resourceType], resourceData)
		}
	}
	
	input["resources"] = resourcesByType
	return input
} 