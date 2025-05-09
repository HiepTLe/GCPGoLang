package misconfig_scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	projectID     string
	reportFormat  string
	outputPath    string
	verbose       bool
	wizClientID   string
	wizClientSecret string
	integrateWiz  bool
	scanType      string
)

// WizAuthResponse represents the response from Wiz authentication
type WizAuthResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// WizVulnerability represents a vulnerability found by Wiz
type WizVulnerability struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Severity      string    `json:"severity"`
	ResourceName  string    `json:"resourceName"`
	ResourceType  string    `json:"resourceType"`
	FirstSeen     time.Time `json:"firstSeen"`
	Status        string    `json:"status"`
	Remediation   string    `json:"remediation"`
	CVE           string    `json:"cve,omitempty"`
}

// Misconfiguration represents a detected GCP configuration issue
type Misconfiguration struct {
	ResourceType  string    `json:"resource_type"`
	ResourceName  string    `json:"resource_name"`
	ResourceID    string    `json:"resource_id"`
	Issue         string    `json:"issue"`
	Severity      string    `json:"severity"`
	Recommendation string   `json:"recommendation"`
	Timestamp     time.Time `json:"timestamp"`
	Category      string    `json:"category"`
}

// ScanResult represents the full scan result
type ScanResult struct {
	ProjectID        string              `json:"project_id"`
	ScanTime         time.Time           `json:"scan_time"`
	Misconfigurations []Misconfiguration  `json:"misconfigurations"`
	WizVulnerabilities []WizVulnerability `json:"wiz_vulnerabilities,omitempty"`
	TotalIssues      int                 `json:"total_issues"`
	SeverityCounts   map[string]int      `json:"severity_counts"`
}

// GetCommand returns the misconfig-scanner command
func GetCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "misconfig-scanner",
		Short: "GCP Misconfiguration Scanner",
		Long: `Scan GCP resources for security misconfigurations and integrate with Wiz
for comprehensive vulnerability management.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create context
			ctx := context.Background()

			if verbose {
				fmt.Printf("Scanning project %s for misconfigurations\n", projectID)
				if integrateWiz {
					fmt.Println("Wiz integration enabled for vulnerability data")
				}
			}

			// Initialize scan result
			result := &ScanResult{
				ProjectID:        projectID,
				ScanTime:         time.Now(),
				Misconfigurations: []Misconfiguration{},
				SeverityCounts:   map[string]int{},
			}

			// Scan for GCP misconfigurations
			if err := scanGCPMisconfigurations(ctx, result); err != nil {
				fmt.Printf("Error scanning GCP resources: %v\n", err)
				os.Exit(1)
			}

			// If Wiz integration is enabled, get vulnerability data
			if integrateWiz && wizClientID != "" && wizClientSecret != "" {
				if err := getWizVulnerabilities(ctx, result); err != nil {
					fmt.Printf("Error getting Wiz vulnerability data: %v\n", err)
					// Continue with GCP results only
				}
			}

			// Count total issues and by severity
			countIssues(result)

			// Output results
			if err := outputResults(result); err != nil {
				fmt.Printf("Error outputting results: %v\n", err)
				os.Exit(1)
			}

			if verbose {
				fmt.Printf("Scan complete. Found %d total issues.\n", result.TotalIssues)
				for severity, count := range result.SeverityCounts {
					fmt.Printf("  %s: %d\n", severity, count)
				}
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.misconfig-scanner.yaml)")
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP project ID")
	rootCmd.PersistentFlags().StringVar(&reportFormat, "report-format", "text", "Output format (text, json, csv)")
	rootCmd.PersistentFlags().StringVar(&outputPath, "output", "", "Output file path (default is stdout)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&integrateWiz, "wiz", false, "Enable Wiz integration for vulnerability data")
	rootCmd.PersistentFlags().StringVar(&wizClientID, "wiz-client-id", "", "Wiz API Client ID")
	rootCmd.PersistentFlags().StringVar(&wizClientSecret, "wiz-client-secret", "", "Wiz API Client Secret")
	rootCmd.PersistentFlags().StringVar(&scanType, "scan-type", "all", "Scan type (all, storage, compute, network, iam)")

	rootCmd.MarkPersistentFlagRequired("project")

	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("report-format", rootCmd.PersistentFlags().Lookup("report-format"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("wiz", rootCmd.PersistentFlags().Lookup("wiz"))
	viper.BindPFlag("wiz-client-id", rootCmd.PersistentFlags().Lookup("wiz-client-id"))
	viper.BindPFlag("wiz-client-secret", rootCmd.PersistentFlags().Lookup("wiz-client-secret"))
	viper.BindPFlag("scan-type", rootCmd.PersistentFlags().Lookup("scan-type"))

	return rootCmd
}

// scanGCPMisconfigurations scans the GCP project for misconfigurations
func scanGCPMisconfigurations(ctx context.Context, result *ScanResult) error {
	// This would be implemented with actual GCP API calls
	// For now, we'll add some example misconfigurations
	
	// Example storage misconfigurations
	if scanType == "all" || scanType == "storage" {
		result.Misconfigurations = append(result.Misconfigurations, Misconfiguration{
			ResourceType:    "storage.googleapis.com/Bucket",
			ResourceName:    "example-bucket",
			ResourceID:      fmt.Sprintf("projects/%s/buckets/example-bucket", projectID),
			Issue:           "Public access enabled",
			Severity:        "HIGH",
			Recommendation:  "Configure uniform bucket-level access and remove public access",
			Timestamp:       time.Now(),
			Category:        "Storage",
		})
		
		result.Misconfigurations = append(result.Misconfigurations, Misconfiguration{
			ResourceType:    "storage.googleapis.com/Bucket",
			ResourceName:    "logs-bucket",
			ResourceID:      fmt.Sprintf("projects/%s/buckets/logs-bucket", projectID),
			Issue:           "Encryption not configured",
			Severity:        "MEDIUM",
			Recommendation:  "Enable CMEK encryption for sensitive data",
			Timestamp:       time.Now(),
			Category:        "Storage",
		})
	}
	
	// Example compute misconfigurations
	if scanType == "all" || scanType == "compute" {
		result.Misconfigurations = append(result.Misconfigurations, Misconfiguration{
			ResourceType:    "compute.googleapis.com/Instance",
			ResourceName:    "instance-1",
			ResourceID:      fmt.Sprintf("projects/%s/zones/us-central1-a/instances/instance-1", projectID),
			Issue:           "Instance has public IP with open SSH port",
			Severity:        "HIGH",
			Recommendation:  "Use IAP for SSH access instead of open firewall rules",
			Timestamp:       time.Now(),
			Category:        "Compute",
		})
	}
	
	// Example network misconfigurations
	if scanType == "all" || scanType == "network" {
		result.Misconfigurations = append(result.Misconfigurations, Misconfiguration{
			ResourceType:    "compute.googleapis.com/Firewall",
			ResourceName:    "default-allow-all",
			ResourceID:      fmt.Sprintf("projects/%s/global/firewalls/default-allow-all", projectID),
			Issue:           "Overly permissive firewall rule (0.0.0.0/0)",
			Severity:        "CRITICAL",
			Recommendation:  "Restrict firewall rules to specific IP ranges",
			Timestamp:       time.Now(),
			Category:        "Network",
		})
	}
	
	// Example IAM misconfigurations
	if scanType == "all" || scanType == "iam" {
		result.Misconfigurations = append(result.Misconfigurations, Misconfiguration{
			ResourceType:    "iam.googleapis.com/ServiceAccount",
			ResourceName:    "service-account-1",
			ResourceID:      fmt.Sprintf("projects/%s/serviceAccounts/service-account-1@%s.iam.gserviceaccount.com", projectID, projectID),
			Issue:           "Service account has owner role",
			Severity:        "HIGH",
			Recommendation:  "Follow principle of least privilege and assign more specific roles",
			Timestamp:       time.Now(),
			Category:        "IAM",
		})
	}

	return nil
}

// getWizVulnerabilities fetches vulnerability data from Wiz API
func getWizVulnerabilities(ctx context.Context, result *ScanResult) error {
	// Authenticate with Wiz API
	token, err := authenticateWiz(wizClientID, wizClientSecret)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Wiz: %w", err)
	}
	
	// In a real implementation, we would query the Wiz API with GraphQL
	// using the authentication token
	if verbose {
		fmt.Printf("Using Wiz token: %s...\n", token[:10])
	}
	
	// For this example, we'll add sample vulnerabilities
	result.WizVulnerabilities = append(result.WizVulnerabilities, WizVulnerability{
		ID:           "wiz-vuln-1",
		Name:         "CVE-2023-1234",
		Description:  "Critical vulnerability in container image",
		Severity:     "CRITICAL",
		ResourceName: "frontend-app",
		ResourceType: "Container",
		FirstSeen:    time.Now().Add(-48 * time.Hour),
		Status:       "OPEN",
		Remediation:  "Update to latest version",
		CVE:          "CVE-2023-1234",
	})
	
	result.WizVulnerabilities = append(result.WizVulnerabilities, WizVulnerability{
		ID:           "wiz-vuln-2",
		Name:         "Outdated TLS Configuration",
		Description:  "Load balancer using outdated TLS configuration",
		Severity:     "MEDIUM",
		ResourceName: "frontend-lb",
		ResourceType: "LoadBalancer",
		FirstSeen:    time.Now().Add(-72 * time.Hour),
		Status:       "OPEN",
		Remediation:  "Update TLS configuration to use TLS 1.2+",
	})
	
	if verbose {
		fmt.Printf("Retrieved %d vulnerabilities from Wiz\n", len(result.WizVulnerabilities))
	}
	
	return nil
}

// authenticateWiz authenticates with the Wiz API and returns a token
func authenticateWiz(clientID, clientSecret string) (string, error) {
	// In a real implementation, we would call the Wiz authentication API
	// For this example, we'll just return a dummy token
	
	if verbose {
		fmt.Println("Authenticating with Wiz API...")
	}
	
	// Create a JWT token (this is just an example, not how Wiz actually works)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": clientID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}
	
	return tokenString, nil
}

// countIssues counts the total issues and issues by severity
func countIssues(result *ScanResult) {
	// Reset counts
	result.TotalIssues = 0
	result.SeverityCounts = map[string]int{
		"CRITICAL": 0,
		"HIGH":     0,
		"MEDIUM":   0,
		"LOW":      0,
	}
	
	// Count GCP misconfigurations
	for _, misc := range result.Misconfigurations {
		result.TotalIssues++
		result.SeverityCounts[misc.Severity]++
	}
	
	// Count Wiz vulnerabilities
	for _, vuln := range result.WizVulnerabilities {
		result.TotalIssues++
		result.SeverityCounts[vuln.Severity]++
	}
}

// outputResults outputs the scan results in the specified format
func outputResults(result *ScanResult) error {
	var output []byte
	var err error
	
	switch strings.ToLower(reportFormat) {
	case "json":
		output, err = json.MarshalIndent(result, "", "  ")
	case "csv":
		// In a real implementation, we would convert to CSV
		output = []byte("CSV output not implemented")
	default:
		// Text format
		output = formatTextOutput(result)
	}
	
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	
	if outputPath == "" {
		// Output to stdout
		fmt.Println(string(output))
	} else {
		// Output to file
		if err := ioutil.WriteFile(outputPath, output, 0644); err != nil {
			return fmt.Errorf("failed to write output to file: %w", err)
		}
	}
	
	return nil
}

// formatTextOutput formats the scan results as text
func formatTextOutput(result *ScanResult) []byte {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("GCP Misconfiguration Scan Results\n"))
	sb.WriteString(fmt.Sprintf("Project: %s\n", result.ProjectID))
	sb.WriteString(fmt.Sprintf("Scan Time: %s\n\n", result.ScanTime.Format(time.RFC3339)))
	
	sb.WriteString(fmt.Sprintf("Total Issues Found: %d\n", result.TotalIssues))
	sb.WriteString(fmt.Sprintf("  CRITICAL: %d\n", result.SeverityCounts["CRITICAL"]))
	sb.WriteString(fmt.Sprintf("  HIGH: %d\n", result.SeverityCounts["HIGH"]))
	sb.WriteString(fmt.Sprintf("  MEDIUM: %d\n", result.SeverityCounts["MEDIUM"]))
	sb.WriteString(fmt.Sprintf("  LOW: %d\n\n", result.SeverityCounts["LOW"]))
	
	sb.WriteString("GCP Misconfigurations:\n")
	for i, misc := range result.Misconfigurations {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, misc.Severity, misc.ResourceName, misc.Issue))
		sb.WriteString(fmt.Sprintf("   Resource: %s\n", misc.ResourceType))
		sb.WriteString(fmt.Sprintf("   Recommendation: %s\n", misc.Recommendation))
		sb.WriteString("\n")
	}
	
	if len(result.WizVulnerabilities) > 0 {
		sb.WriteString("Wiz Vulnerabilities:\n")
		for i, vuln := range result.WizVulnerabilities {
			sb.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, vuln.Severity, vuln.ResourceName, vuln.Name))
			sb.WriteString(fmt.Sprintf("   Description: %s\n", vuln.Description))
			sb.WriteString(fmt.Sprintf("   First Seen: %s\n", vuln.FirstSeen.Format(time.RFC3339)))
			sb.WriteString(fmt.Sprintf("   Remediation: %s\n", vuln.Remediation))
			if vuln.CVE != "" {
				sb.WriteString(fmt.Sprintf("   CVE: %s\n", vuln.CVE))
			}
			sb.WriteString("\n")
		}
	}
	
	return []byte(sb.String())
} 