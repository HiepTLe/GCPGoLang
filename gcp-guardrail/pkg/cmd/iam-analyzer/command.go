package iam_analyzer

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hieptle/gcp-guardrail/pkg/gcp/iam"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	projectID    string
	reportFormat string
	outputPath   string
	verbose      bool
	riskLevel    string
)

// GetCommand returns the iam-analyzer command
func GetCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "iam-analyzer",
		Short: "GCP IAM policy analyzer",
		Long: `Analyze GCP IAM policies to identify overly permissive permissions,
policy violations, and security risks.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create context
			// ctx := context.Background()

			if verbose {
				fmt.Printf("Analyzing IAM policies for project %s\n", projectID)
				fmt.Printf("Risk level filter: %s\n", riskLevel)
				fmt.Printf("Analysis started at: %s\n", time.Now().Format(time.RFC3339))
			}
			
			// Convert risk level to integer
			riskInt, err := strconv.Atoi(riskLevel)
			if err != nil {
				fmt.Printf("Warning: Invalid risk level '%s', using default (3)\n", riskLevel)
				riskInt = 3 // Default if not a number
			}
			
			// In a real implementation, we would initialize the analyzer and use it
			// ctx := context.Background()
			// analyzer, err := iam.NewAnalyzer(ctx, projectID)
			// if err != nil {
			//    fmt.Printf("Error: Failed to create IAM analyzer: %v\n", err)
			//    os.Exit(1)
			// }
			
			// For now, we'll create a sample analysis with test data
			// In the future, this would call analyzer.AnalyzeProject()
			analysis := &iam.Analysis{
				ProjectID: projectID,
				Timestamp: time.Now(),
				Issues: []iam.Issue{
					{
						Severity:    "CRITICAL",
						Description: "User account has Owner role at organization level",
						Principal:   "user:admin@example.com",
						Role:        "roles/owner",
						Mitigation:  "Remove Owner role and grant more specific roles",
					},
					{
						Severity:    "HIGH",
						Description: "Service account has broad permissions",
						Principal:   "serviceAccount:sa@project.iam.gserviceaccount.com",
						Role:        "roles/editor",
						Mitigation:  "Grant only required permissions to service account",
					},
					{
						Severity:    "MEDIUM",
						Description: "Group has compute admin permissions",
						Principal:   "group:engineers@example.com",
						Role:        "roles/compute.admin",
						Mitigation:  "Limit compute admin access to specific principals",
					},
					{
						Severity:    "LOW",
						Description: "User has viewer permissions across multiple projects",
						Principal:   "user:viewer@example.com",
						Role:        "roles/viewer",
						Mitigation:  "Review necessity for cross-project access",
					},
				},
				RoleAssignments: []iam.RoleAssignment{
					{
						Principal: "user:admin@example.com",
						Role:      "roles/owner",
						Scope:     "organization/123456789",
					},
					{
						Principal: "serviceAccount:sa@project.iam.gserviceaccount.com",
						Role:      "roles/editor",
						Scope:     "project/" + projectID,
					},
				},
			}
			
			// Filter issues based on risk level
			var filteredIssues []iam.Issue
			for _, issue := range analysis.Issues {
				// Convert severity to risk level (simplified mapping)
				var issueRisk int
				switch issue.Severity {
				case "CRITICAL":
					issueRisk = 5
				case "HIGH":
					issueRisk = 4
				case "MEDIUM":
					issueRisk = 3
				case "LOW":
					issueRisk = 2
				default:
					issueRisk = 1
				}
				
				if issueRisk >= riskInt {
					filteredIssues = append(filteredIssues, issue)
				}
			}
			analysis.Issues = filteredIssues

			// Create report from analysis
			report := iam.NewReport(analysis)
			
			// Determine report format
			var format iam.ReportFormat
			switch reportFormat {
			case "json":
				format = iam.JSONFormat
			case "csv":
				format = iam.CSVFormat
			default:
				format = iam.TextFormat
			}
			
			// Output the report
			if err := iam.WriteReportToFile(outputPath, report, format); err != nil {
				fmt.Printf("Error: Failed to write report: %v\n", err)
				os.Exit(1)
			}
			
			if verbose {
				fmt.Printf("Analysis completed. Found %d policy violations.\n", len(analysis.Issues))
				if outputPath != "" {
					fmt.Printf("Report written to %s\n", outputPath)
				}
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iam-analyzer.yaml)")
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP project ID")
	rootCmd.PersistentFlags().StringVar(&reportFormat, "report-format", "text", "Output format (text, json, csv)")
	rootCmd.PersistentFlags().StringVar(&outputPath, "output", "", "Output file path (default is stdout)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&riskLevel, "risk-level", "3", "Minimum risk level to report (1-5)")

	rootCmd.MarkPersistentFlagRequired("project")

	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("report-format", rootCmd.PersistentFlags().Lookup("report-format"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("risk-level", rootCmd.PersistentFlags().Lookup("risk-level"))

	return rootCmd
} 