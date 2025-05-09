package iam

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// ReportFormat defines the format of the IAM analysis report
type ReportFormat string

const (
	// TextFormat outputs the report in a human-readable text format
	TextFormat ReportFormat = "text"
	// JSONFormat outputs the report in JSON format
	JSONFormat ReportFormat = "json"
	// CSVFormat outputs the report in CSV format
	CSVFormat ReportFormat = "csv"
)

// Report represents an IAM analysis report
type Report struct {
	ProjectID         string             `json:"project_id"`
	GeneratedAt       time.Time          `json:"generated_at"`
	RoleAssignments   []RoleAssignment   `json:"role_assignments,omitempty"`
	Issues            []Issue            `json:"issues"`
	UnusedPermissions []UnusedPermission `json:"unused_permissions,omitempty"`
	Stats             struct {
		TotalIssues       int `json:"total_issues"`
		CriticalIssues    int `json:"critical_issues"`
		HighIssues        int `json:"high_issues"`
		MediumIssues      int `json:"medium_issues"`
		LowIssues         int `json:"low_issues"`
		TotalRoles        int `json:"total_roles"`
		TotalPrincipals   int `json:"total_principals"`
		UnusedPermissions int `json:"unused_permissions"`
	} `json:"stats"`
}

// NewReport creates a new IAM analysis report from an Analysis result
func NewReport(analysis *Analysis) *Report {
	report := &Report{
		ProjectID:         analysis.ProjectID,
		GeneratedAt:       time.Now(),
		RoleAssignments:   analysis.RoleAssignments,
		Issues:            analysis.Issues,
		UnusedPermissions: analysis.UnusedPermissions,
	}

	// Calculate stats
	report.Stats.TotalIssues = len(analysis.Issues)
	report.Stats.TotalRoles = len(analysis.RoleAssignments)
	report.Stats.UnusedPermissions = len(analysis.UnusedPermissions)

	// Count unique principals
	principals := make(map[string]bool)
	for _, ra := range analysis.RoleAssignments {
		principals[ra.Principal] = true
	}
	report.Stats.TotalPrincipals = len(principals)

	// Count issues by severity
	for _, issue := range analysis.Issues {
		switch issue.Severity {
		case "CRITICAL":
			report.Stats.CriticalIssues++
		case "HIGH":
			report.Stats.HighIssues++
		case "MEDIUM":
			report.Stats.MediumIssues++
		case "LOW":
			report.Stats.LowIssues++
		}
	}

	return report
}

// WriteReport writes the report to the specified writer in the specified format
func WriteReport(w io.Writer, report *Report, format ReportFormat) error {
	switch format {
	case TextFormat:
		return writeTextReport(w, report)
	case JSONFormat:
		return writeJSONReport(w, report)
	case CSVFormat:
		return writeCSVReport(w, report)
	default:
		return fmt.Errorf("unsupported report format: %s", format)
	}
}

// writeTextReport writes the report in a human-readable text format
func writeTextReport(w io.Writer, report *Report) error {
	// Write header
	fmt.Fprintf(w, "# IAM POLICY ANALYSIS REPORT\n")
	fmt.Fprintf(w, "Project: %s\n", report.ProjectID)
	fmt.Fprintf(w, "Generated: %s\n\n", report.GeneratedAt.Format(time.RFC1123))

	// Write stats
	fmt.Fprintf(w, "## SUMMARY\n")
	fmt.Fprintf(w, "Total issues: %d\n", report.Stats.TotalIssues)
	fmt.Fprintf(w, "  Critical: %d\n", report.Stats.CriticalIssues)
	fmt.Fprintf(w, "  High: %d\n", report.Stats.HighIssues)
	fmt.Fprintf(w, "  Medium: %d\n", report.Stats.MediumIssues)
	fmt.Fprintf(w, "  Low: %d\n", report.Stats.LowIssues)
	fmt.Fprintf(w, "Total roles analyzed: %d\n", report.Stats.TotalRoles)
	fmt.Fprintf(w, "Total principals: %d\n", report.Stats.TotalPrincipals)
	fmt.Fprintf(w, "Unused permissions: %d\n\n", report.Stats.UnusedPermissions)

	// Write security issues
	fmt.Fprintf(w, "## SECURITY ISSUES\n")
	for i, issue := range report.Issues {
		fmt.Fprintf(w, "%d. [%s] %s\n", i+1, issue.Severity, issue.Description)
		fmt.Fprintf(w, "   Principal: %s\n", issue.Principal)
		fmt.Fprintf(w, "   Role: %s\n", issue.Role)
		fmt.Fprintf(w, "   Mitigation: %s\n\n", issue.Mitigation)
	}

	// Write role assignments if present
	if len(report.RoleAssignments) > 0 {
		fmt.Fprintf(w, "## ROLE ASSIGNMENTS\n")
		for i, ra := range report.RoleAssignments {
			fmt.Fprintf(w, "%d. Principal: %s\n", i+1, ra.Principal)
			fmt.Fprintf(w, "   Role: %s\n", ra.Role)
			fmt.Fprintf(w, "   Scope: %s\n", ra.Scope)
			
			if len(ra.EffectiveAccess) > 0 {
				fmt.Fprintf(w, "   Effective Access:\n")
				for _, access := range ra.EffectiveAccess {
					fmt.Fprintf(w, "     - %s\n", access)
				}
			}
			fmt.Fprintf(w, "\n")
		}
	}

	// Write unused permissions if present
	if len(report.UnusedPermissions) > 0 {
		fmt.Fprintf(w, "## UNUSED PERMISSIONS\n")
		for i, up := range report.UnusedPermissions {
			fmt.Fprintf(w, "%d. Principal: %s\n", i+1, up.Principal)
			fmt.Fprintf(w, "   Role: %s\n", up.Role)
			fmt.Fprintf(w, "   Permission: %s\n", up.Permission)
			fmt.Fprintf(w, "   Last Used: %s\n", formatLastUsed(up.LastUsed))
			fmt.Fprintf(w, "   Recommended: %s\n\n", up.Recommended)
		}
	}

	return nil
}

// writeJSONReport writes the report in JSON format
func writeJSONReport(w io.Writer, report *Report) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// writeCSVReport writes the report in CSV format
func writeCSVReport(w io.Writer, report *Report) error {
	csvWriter := csv.NewWriter(w)
	
	// Write issues header
	headers := []string{
		"Severity", "Description", "Principal", "Role", "Mitigation",
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}
	
	// Write issues
	for _, issue := range report.Issues {
		row := []string{
			issue.Severity,
			issue.Description,
			issue.Principal,
			issue.Role,
			issue.Mitigation,
		}
		
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}
	
	csvWriter.Flush()
	return csvWriter.Error()
}

// formatLastUsed formats the last used time for display
func formatLastUsed(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return t.Format(time.RFC1123)
}

// WriteReportToFile writes the report to a file
func WriteReportToFile(filename string, report *Report, format ReportFormat) error {
	var file *os.File
	var err error
	
	if filename == "" {
		file = os.Stdout
	} else {
		file, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
	}
	
	return WriteReport(file, report, format)
} 