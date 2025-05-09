package sa_tracker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hieptle/gcp-guardrail/pkg/gcp/sa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	projectID    string
	reportFormat string
	outputPath   string
	daysBack     int
	verbose      bool
)

// GetCommand returns the sa-tracker command
func GetCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sa-tracker",
		Short: "GCP Service Account usage tracker",
		Long: `Track and analyze GCP Service Account usage patterns to identify unused
accounts, over-permissioned accounts, and anomalous behavior.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create context
			ctx := context.Background()

			// Create a service account tracker
			tracker, err := sa.NewTracker(ctx, projectID)
			if err != nil {
				fmt.Printf("Failed to create service account tracker: %v\n", err)
				os.Exit(1)
			}
			defer tracker.Close()

			// Convert daysBack to a duration for lookback period
			lookbackPeriod := time.Duration(daysBack) * 24 * time.Hour
			
			if verbose {
				fmt.Printf("Analyzing service account usage for project %s\n", projectID)
				fmt.Printf("Looking back %d days of activity (%s)\n", daysBack, lookbackPeriod)
			}

			// Run the service account analysis
			serviceAccounts, err := tracker.AnalyzeUsage(lookbackPeriod)
			if err != nil {
				fmt.Printf("Failed to analyze service account usage: %v\n", err)
				os.Exit(1)
			}

			// Create a report
			report := sa.NewReport(projectID, lookbackPeriod, serviceAccounts)
			
			// Determine report format
			var format sa.ReportFormat
			switch reportFormat {
			case "json":
				format = sa.JSONFormat
			case "csv":
				format = sa.CSVFormat
			default:
				format = sa.TextFormat
			}
			
			// Output the report
			if err := sa.WriteReportToFile(outputPath, report, format); err != nil {
				fmt.Printf("Failed to write report: %v\n", err)
				os.Exit(1)
			}

			// Print unused accounts count if verbose
			if verbose {
				fmt.Printf("Found %d total service accounts\n", report.Stats.TotalAccounts)
				fmt.Printf("Found %d unused service accounts (%s%%)\n", 
					report.Stats.UnusedAccounts, 
					strconv.FormatFloat(float64(report.Stats.UnusedAccounts)/float64(report.Stats.TotalAccounts)*100, 'f', 1, 64))
				fmt.Printf("Found %d over-privileged service accounts\n", report.Stats.OverPrivAccounts)
				
				if outputPath != "" {
					fmt.Printf("Report written to %s in %s format\n", outputPath, reportFormat)
				}
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sa-tracker.yaml)")
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP project ID")
	rootCmd.PersistentFlags().StringVar(&reportFormat, "report-format", "text", "Output format (text, json, csv)")
	rootCmd.PersistentFlags().StringVar(&outputPath, "output", "", "Output file path (default is stdout)")
	rootCmd.PersistentFlags().IntVar(&daysBack, "days", 30, "Number of days to look back for service account activity")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")

	rootCmd.MarkPersistentFlagRequired("project")

	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("report-format", rootCmd.PersistentFlags().Lookup("report-format"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("days", rootCmd.PersistentFlags().Lookup("days"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	
	return rootCmd
} 