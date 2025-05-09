package main

import (
	"fmt"
	"os"

	iam "github.com/hieptle/gcp-guardrail/pkg/cmd/iam-analyzer"
	log "github.com/hieptle/gcp-guardrail/pkg/cmd/log-watcher"
	misconfig "github.com/hieptle/gcp-guardrail/pkg/cmd/misconfig-scanner"
	sa "github.com/hieptle/gcp-guardrail/pkg/cmd/sa-tracker"
	tf "github.com/hieptle/gcp-guardrail/pkg/cmd/tf-validator"
	"github.com/spf13/cobra"
)

var (
	projectID string
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "gcpgolang",
	Short: "GCPGoLang - GCP Security Suite",
	Long: `GCPGoLang is a comprehensive security suite for Google Cloud Platform
using Golang and Rego policies. It provides security analysis, monitoring, 
and policy enforcement for your GCP resources.

Project information:
- Project name: GCPGoLang
- Project number: 652769711122
- Project ID: gcpgolang`,
	Run: func(cmd *cobra.Command, args []string) {
		if projectID == "" {
			projectID = "gcpgolang" // default to our project ID
		}
		fmt.Printf("GCPGoLang Security Suite\n")
		fmt.Printf("Project ID: %s\n", projectID)
		fmt.Printf("\nUse --help to see available commands\n")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP project ID (default: gcpgolang)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")

	// Add all the subcommands from gcp-guardrail
	rootCmd.AddCommand(iam.GetCommand())
	rootCmd.AddCommand(sa.GetCommand())
	rootCmd.AddCommand(log.GetCommand())
	rootCmd.AddCommand(tf.GetCommand())
	rootCmd.AddCommand(misconfig.GetCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 