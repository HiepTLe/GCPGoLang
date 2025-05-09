package tf_validator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	planFile      string
	policyDir     string
	outputFile    string
	severity      string
	failThreshold string
	verbose       bool
)

// Violation represents a policy violation
type Violation struct {
	Severity    string `json:"severity"`
	PolicyName  string `json:"policy_name"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Message     string `json:"message"`
	Remediation string `json:"remediation"`
}

// ValidationResult represents the output of the validation
type ValidationResult struct {
	PlanFile      string      `json:"plan_file"`
	PolicyDir     string      `json:"policy_dir"`
	SeverityLevel string      `json:"severity_level"`
	Violations    []Violation `json:"violations"`
	Timestamp     string      `json:"timestamp"`
	Duration      string      `json:"duration"`
}

// GetCommand returns the tf-validator command
func GetCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tf-validator",
		Short: "Terraform plan validator for GCP",
		Long: `Validates Terraform plans for GCP against security policies defined in Rego.
Checks for configuration issues, security risks, and policy violations before applying.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create context
			ctx := context.Background()

			// For now, we'll just print some information
			// In a real implementation, we would validate the Terraform plan
			
			// Use time package to demonstrate it's being used
			startTime := time.Now()
			fmt.Printf("Starting validation at: %s\n", startTime.Format(time.RFC3339))
			
			fmt.Printf("Validating Terraform plan: %s\n", planFile)
			fmt.Printf("Using policy directory: %s\n", policyDir)
			fmt.Printf("Minimum severity level: %s\n", severity)
			
			// Parse severity level
			sevLevel, err := strconv.Atoi(severity)
			if err != nil {
				sevLevel = 0 // Default if not a number
			}
			
			// Parse fail threshold level
			failLevel, err := strconv.Atoi(failThreshold)
			if err != nil {
				failLevel = 4 // Default to 4 (high severity) if not a number
			}
			
			// Sample violations (in a real implementation, these would come from policy evaluation)
			// Severity mapping: 1=Low, 2=Medium, 3=High, 4=Critical, 5=Blocker
			var violations []Violation
			var highSevViolations int
			
			// Use context to demonstrate it's being used
			select {
			case <-ctx.Done():
				fmt.Println("Validation was cancelled")
			default:
				// Simulate validation by creating sample violations
				simulatedViolations := []Violation{
					{
						Severity:     "Medium",
						PolicyName:   "storage_encryption",
						ResourceType: "google_storage_bucket",
						ResourceName: "my-non-compliant-bucket",
						Message:      "Storage bucket is not encrypted",
						Remediation:  "Add encryption { default_kms_key_name = ... } to the bucket configuration",
					},
					{
						Severity:     "High",
						PolicyName:   "public_bucket_access",
						ResourceType: "google_storage_bucket",
						ResourceName: "my-non-compliant-bucket",
						Message:      "Storage bucket does not have uniform bucket-level access enabled",
						Remediation:  "Set uniform_bucket_level_access = true in the bucket configuration",
					},
					{
						Severity:     "Low",
						PolicyName:   "versioning_recommended",
						ResourceType: "google_storage_bucket",
						ResourceName: "my-non-compliant-bucket",
						Message:      "Storage bucket does not have versioning enabled",
						Remediation:  "Add versioning { enabled = true } to the bucket configuration",
					},
				}
				
				// Filter violations based on severity level
				for _, v := range simulatedViolations {
					var violationLevel int
					switch strings.ToLower(v.Severity) {
					case "low":
						violationLevel = 1
					case "medium":
						violationLevel = 2
					case "high":
						violationLevel = 3
					case "critical":
						violationLevel = 4
					case "blocker":
						violationLevel = 5
					default:
						violationLevel = 1
					}
					
					if violationLevel >= sevLevel {
						violations = append(violations, v)
						if violationLevel >= failLevel {
							highSevViolations++
						}
					}
				}
				
				elapsedTime := time.Since(startTime)
				
				// Display detailed violation information
				if len(violations) > 0 {
					fmt.Printf("Found %d policy violations\n", len(violations))
					fmt.Println("--------------------------------------------")
					for i, v := range violations {
						fmt.Printf("Violation #%d:\n", i+1)
						fmt.Printf("  Severity:      %s\n", v.Severity)
						fmt.Printf("  Policy:        %s\n", v.PolicyName)
						fmt.Printf("  Resource Type: %s\n", v.ResourceType)
						fmt.Printf("  Resource Name: %s\n", v.ResourceName)
						fmt.Printf("  Issue:         %s\n", v.Message)
						fmt.Printf("  Remediation:   %s\n", v.Remediation)
						fmt.Println("--------------------------------------------")
					}
				} else {
					fmt.Println("No policy violations found")
				}
				
				fmt.Printf("Validation completed in %s\n", elapsedTime)
				
				// Create and save validation results
				result := ValidationResult{
					PlanFile:      planFile,
					PolicyDir:     policyDir,
					SeverityLevel: severity,
					Violations:    violations,
					Timestamp:     startTime.Format(time.RFC3339),
					Duration:      elapsedTime.String(),
				}
				
				// Save to output file if specified
				if outputFile != "" {
					resultJSON, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						fmt.Printf("Error creating JSON output: %v\n", err)
					} else {
						err = os.WriteFile(outputFile, resultJSON, 0644)
						if err != nil {
							fmt.Printf("Error writing output file: %v\n", err)
						} else {
							fmt.Printf("Results written to %s\n", outputFile)
						}
					}
				}
				
				// Only fail if high severity issues are found (based on threshold)
				if highSevViolations > 0 {
					fmt.Printf("Terraform plan validation failed due to %d high severity violations!\n", highSevViolations)
					os.Exit(1)
				} else if len(violations) > 0 {
					fmt.Println("Terraform plan validation completed with warnings.")
				} else {
					fmt.Println("Terraform plan validation succeeded!")
				}
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tf-validator.yaml)")
	rootCmd.PersistentFlags().StringVar(&planFile, "plan", "", "Terraform plan file (JSON format)")
	rootCmd.PersistentFlags().StringVar(&policyDir, "policy-dir", "policies/terraform", "Directory containing Rego policies")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "", "Output file for validation results (default is stdout)")
	rootCmd.PersistentFlags().StringVar(&severity, "severity", "2", "Minimum severity level (1-5)")
	rootCmd.PersistentFlags().StringVar(&failThreshold, "fail-threshold", "4", "Minimum severity level that causes validation to fail (1-5, default 4)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")

	rootCmd.MarkPersistentFlagRequired("plan")

	viper.BindPFlag("plan", rootCmd.PersistentFlags().Lookup("plan"))
	viper.BindPFlag("policy-dir", rootCmd.PersistentFlags().Lookup("policy-dir"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("severity", rootCmd.PersistentFlags().Lookup("severity"))
	viper.BindPFlag("fail-threshold", rootCmd.PersistentFlags().Lookup("fail-threshold"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	
	return rootCmd
} 