package sa

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/api/logging/v2"
)

// ServiceAccount represents a GCP service account with usage information
type ServiceAccount struct {
	Email         string    `json:"email"`
	DisplayName   string    `json:"display_name"`
	LastUsed      time.Time `json:"last_used"`
	KeyCount      int       `json:"key_count"`
	IsUsed        bool      `json:"is_used"`
	IsOverPriv    bool      `json:"is_over_privileged"`
	Roles         []string  `json:"roles"`
	Created       time.Time `json:"created"`
	ActivityCount int       `json:"activity_count"`
}

// Tracker analyzes service account usage in a GCP project
type Tracker struct {
	projectID         string
	loggingService    *logging.Service
	ctx               context.Context
}

// NewTracker creates a new service account tracker for a GCP project
func NewTracker(ctx context.Context, projectID string) (*Tracker, error) {
	// Create logging service for activity tracking
	loggingService, err := logging.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create logging service: %w", err)
	}

	return &Tracker{
		projectID:      projectID,
		loggingService: loggingService,
		ctx:            ctx,
	}, nil
}

// AnalyzeUsage analyzes service account usage over a period of time
func (t *Tracker) AnalyzeUsage(lookbackPeriod time.Duration) ([]*ServiceAccount, error) {
	// For this implementation, we'll simulate fetching service accounts
	// In a real implementation, you would use the IAM Admin API
	serviceAccounts := []*ServiceAccount{
		{
			Email:       fmt.Sprintf("sa-1@%s.iam.gserviceaccount.com", t.projectID),
			DisplayName: "Service Account 1",
			Created:     time.Now().Add(-90 * 24 * time.Hour),
		},
		{
			Email:       fmt.Sprintf("sa-2@%s.iam.gserviceaccount.com", t.projectID),
			DisplayName: "Service Account 2",
			Created:     time.Now().Add(-60 * 24 * time.Hour),
		},
		{
			Email:       fmt.Sprintf("sa-3@%s.iam.gserviceaccount.com", t.projectID),
			DisplayName: "Service Account 3",
			Created:     time.Now().Add(-30 * 24 * time.Hour),
		},
	}

	// For each service account, analyze usage
	for _, sa := range serviceAccounts {
		// Simulate key count data
		sa.KeyCount = len(sa.Email) % 3

		// Check activity logs
		lastUsed, activityCount, err := t.checkActivity(sa.Email, lookbackPeriod)
		if err == nil {
			sa.LastUsed = lastUsed
			sa.ActivityCount = activityCount
			sa.IsUsed = !lastUsed.IsZero() && time.Since(lastUsed) < lookbackPeriod
		}

		// Get roles
		roles, err := t.getRoles(sa.Email)
		if err == nil {
			sa.Roles = roles
			// Simple heuristic for over-privileged accounts
			sa.IsOverPriv = len(roles) > 5 && sa.ActivityCount < 10
		}
	}

	return serviceAccounts, nil
}

// checkActivity checks activity logs for a service account
func (t *Tracker) checkActivity(email string, lookbackPeriod time.Duration) (time.Time, int, error) {
	var lastUsed time.Time
	activityCount := 0
	
	endTime := time.Now()
	startTime := endTime.Add(-lookbackPeriod)
	
	// Format timestamps for the filter
	startTimeStr := startTime.Format(time.RFC3339)
	endTimeStr := endTime.Format(time.RFC3339)
	
	// Create filter for this service account's activity
	filter := fmt.Sprintf(`protoPayload.authenticationInfo.principalEmail="%s" AND timestamp >= "%s" AND timestamp <= "%s"`,
		email, startTimeStr, endTimeStr)
	
	// Create the entries list call
	entriesService := t.loggingService.Entries
	listCall := entriesService.List(&logging.ListLogEntriesRequest{
		ResourceNames: []string{"projects/" + t.projectID},
		Filter:        filter,
		OrderBy:       "timestamp desc",
		PageSize:      1000, // Limit to reasonable number
	})
	
	// Collect log entries
	err := listCall.Pages(t.ctx, func(page *logging.ListLogEntriesResponse) error {
		for i, entry := range page.Entries {
			activityCount++
			
			// Record the timestamp of the most recent activity (first entry)
			if i == 0 && lastUsed.IsZero() {
				// Convert timestamp string to time.Time
				if t, err := time.Parse(time.RFC3339, entry.Timestamp); err == nil {
					lastUsed = t
				}
			}
		}
		return nil
	})
	
	if err != nil {
		return lastUsed, activityCount, fmt.Errorf("failed to check activity logs: %w", err)
	}
	
	return lastUsed, activityCount, nil
}

// getRoles gets the roles assigned to a service account
func (t *Tracker) getRoles(email string) ([]string, error) {
	// In a real implementation, you would fetch the IAM policy and extract roles
	// Here we're returning simulated data based on the email to use strconv
	numRoles := len(email) % 5 + 2 // 2-6 roles
	roles := make([]string, numRoles)
	
	for i := 0; i < numRoles; i++ {
		roleNum := strconv.Itoa(i + 1)
		roles[i] = fmt.Sprintf("roles/role%s", roleNum)
	}
	
	return roles, nil
}

// Close closes the tracker and releases resources
func (t *Tracker) Close() error {
	return nil
} 