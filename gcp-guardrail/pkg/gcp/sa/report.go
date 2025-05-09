package sa

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// ReportFormat defines the format of the service account report
type ReportFormat string

const (
	// TextFormat outputs the report in a human-readable text format
	TextFormat ReportFormat = "text"
	// JSONFormat outputs the report in JSON format
	JSONFormat ReportFormat = "json"
	// CSVFormat outputs the report in CSV format
	CSVFormat ReportFormat = "csv"
)

// Report represents a service account usage report
type Report struct {
	ProjectID       string           `json:"project_id"`
	GeneratedAt     time.Time        `json:"generated_at"`
	LookbackPeriod  string           `json:"lookback_period"`
	ServiceAccounts []*ServiceAccount `json:"service_accounts"`
	Stats           struct {
		TotalAccounts      int `json:"total_accounts"`
		UnusedAccounts     int `json:"unused_accounts"`
		OverPrivAccounts   int `json:"over_privileged_accounts"`
		AccountsWithKeys   int `json:"accounts_with_keys"`
		TotalKeys          int `json:"total_keys"`
	} `json:"stats"`
}

// NewReport creates a new service account report
func NewReport(projectID string, lookbackPeriod time.Duration, accounts []*ServiceAccount) *Report {
	report := &Report{
		ProjectID:       projectID,
		GeneratedAt:     time.Now(),
		LookbackPeriod:  lookbackPeriod.String(),
		ServiceAccounts: accounts,
	}

	// Calculate stats
	for _, sa := range accounts {
		report.Stats.TotalAccounts++
		
		if !sa.IsUsed {
			report.Stats.UnusedAccounts++
		}
		
		if sa.IsOverPriv {
			report.Stats.OverPrivAccounts++
		}
		
		if sa.KeyCount > 0 {
			report.Stats.AccountsWithKeys++
			report.Stats.TotalKeys += sa.KeyCount
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
	fmt.Fprintf(w, "# SERVICE ACCOUNT USAGE REPORT\n")
	fmt.Fprintf(w, "Project: %s\n", report.ProjectID)
	fmt.Fprintf(w, "Generated: %s\n", report.GeneratedAt.Format(time.RFC1123))
	fmt.Fprintf(w, "Lookback Period: %s\n\n", report.LookbackPeriod)
	
	// Write stats
	fmt.Fprintf(w, "## SUMMARY\n")
	fmt.Fprintf(w, "Total service accounts: %d\n", report.Stats.TotalAccounts)
	fmt.Fprintf(w, "Unused service accounts: %d\n", report.Stats.UnusedAccounts)
	fmt.Fprintf(w, "Over-privileged service accounts: %d\n", report.Stats.OverPrivAccounts)
	fmt.Fprintf(w, "Service accounts with keys: %d\n", report.Stats.AccountsWithKeys)
	fmt.Fprintf(w, "Total keys: %d\n\n", report.Stats.TotalKeys)
	
	// Write service account details
	fmt.Fprintf(w, "## SERVICE ACCOUNTS\n")
	for i, sa := range report.ServiceAccounts {
		fmt.Fprintf(w, "%d. %s (%s)\n", i+1, sa.Email, sa.DisplayName)
		fmt.Fprintf(w, "   Created: %s\n", sa.Created.Format(time.RFC1123))
		fmt.Fprintf(w, "   Last used: %s\n", formatLastUsed(sa.LastUsed))
		fmt.Fprintf(w, "   Activity count: %d\n", sa.ActivityCount)
		fmt.Fprintf(w, "   Key count: %d\n", sa.KeyCount)
		fmt.Fprintf(w, "   Used: %t\n", sa.IsUsed)
		fmt.Fprintf(w, "   Over-privileged: %t\n", sa.IsOverPriv)
		
		if len(sa.Roles) > 0 {
			fmt.Fprintf(w, "   Roles:\n")
			for _, role := range sa.Roles {
				fmt.Fprintf(w, "     - %s\n", role)
			}
		} else {
			fmt.Fprintf(w, "   Roles: None or unable to retrieve\n")
		}
		
		fmt.Fprintf(w, "\n")
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
	
	// Write header
	headers := []string{
		"Email", "Display Name", "Created", "Last Used", "Activity Count",
		"Key Count", "Is Used", "Is Over-Privileged", "Roles",
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}
	
	// Write service accounts
	for _, sa := range report.ServiceAccounts {
		roles := ""
		for i, role := range sa.Roles {
			if i > 0 {
				roles += "; "
			}
			roles += role
		}
		
		row := []string{
			sa.Email,
			sa.DisplayName,
			sa.Created.Format(time.RFC3339),
			formatLastUsed(sa.LastUsed),
			strconv.Itoa(sa.ActivityCount),
			strconv.Itoa(sa.KeyCount),
			strconv.FormatBool(sa.IsUsed),
			strconv.FormatBool(sa.IsOverPriv),
			roles,
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