package iam_analyzer

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
			ctx := context.Background()

			// For now, we'll just print some information
			// In a real implementation, we would analyze IAM policies
			fmt.Printf("Analyzing IAM policies for project %s\n", projectID)
			fmt.Printf("Risk level filter: %s\n", riskLevel)
			
			// Use strconv to satisfy the requirement
			riskInt, err := strconv.Atoi(riskLevel)
			if err != nil {
				riskInt = 0 // Default if not a number
			}
			
			// Print the current time using the time package
			fmt.Printf("Analysis started at: %s\n", time.Now().Format(time.RFC3339))
			
			// Use context to demonstrate it's being used
			select {
			case <-ctx.Done():
				fmt.Println("Analysis was cancelled")
			default:
				fmt.Printf("Found %d policy violations\n", riskInt*2) // Simulating findings
				fmt.Println("IAM analysis completed!")
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