package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/logging/apiv2/loggingpb"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/logging/v2"
)

// Alert represents a security alert triggered from log analysis
type Alert struct {
	Timestamp   time.Time                 `json:"timestamp"`
	Severity    string                    `json:"severity"`
	Description string                    `json:"description"`
	Resource    string                    `json:"resource"`
	ProjectID   string                    `json:"project_id"`
	LogName     string                    `json:"log_name"`
	Details     map[string]interface{}    `json:"details,omitempty"`
	LogEntry    *loggingpb.LogEntry       `json:"-"`
}

// Monitor watches GCP logs for security events
type Monitor struct {
	projectID    string
	loggingService *logging.Service
	pubsubClient  *pubsub.Client
	alertTopic    *pubsub.Topic
	ctx           context.Context
}

// NewMonitor creates a new GCP logging monitor
func NewMonitor(ctx context.Context, projectID string, alertTopicID string) (*Monitor, error) {
	// Create Logging service
	loggingService, err := logging.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create logging service: %w", err)
	}

	// Create Pub/Sub client for alerts
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	// Get or create the alert topic
	alertTopic := pubsubClient.Topic(alertTopicID)
	exists, err := alertTopic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if !exists {
		alertTopic, err = pubsubClient.CreateTopic(ctx, alertTopicID)
		if err != nil {
			return nil, fmt.Errorf("failed to create alert topic: %w", err)
		}
	}

	return &Monitor{
		projectID:      projectID,
		loggingService: loggingService,
		pubsubClient:   pubsubClient,
		alertTopic:     alertTopic,
		ctx:            ctx,
	}, nil
}

// QueryLogs queries logs and looks for security incidents
func (m *Monitor) QueryLogs(filter string, timeWindow time.Duration) ([]*Alert, error) {
	endTime := time.Now()
	startTime := endTime.Add(-timeWindow)

	// Format timestamps for the filter
	startTimeStr := startTime.Format(time.RFC3339)
	endTimeStr := endTime.Format(time.RFC3339)

	// Append time constraints to the filter
	timeFilter := fmt.Sprintf(`timestamp >= "%s" AND timestamp <= "%s"`, startTimeStr, endTimeStr)
	if filter != "" {
		filter = fmt.Sprintf(`%s AND %s`, filter, timeFilter)
	} else {
		filter = timeFilter
	}

	// Create the entries list call
	entriesService := m.loggingService.Entries
	listCall := entriesService.List(&logging.ListLogEntriesRequest{
		ResourceNames: []string{"projects/" + m.projectID},
		Filter:     filter,
	})

	// Collect log entries and look for security incidents
	var alerts []*Alert
	err := listCall.Pages(m.ctx, func(page *logging.ListLogEntriesResponse) error {
		for _, entry := range page.Entries {
			// Process each log entry
			// Here you would implement specific logic to detect security incidents
			// This is a simplified example
			alert := &Alert{
				Timestamp:   time.Now(),
				Severity:    "INFO",
				Description: fmt.Sprintf("Detected log activity with ID: %s", entry.InsertId),
				Resource:    entry.Resource.Type,
				ProjectID:   m.projectID,
				LogName:     entry.LogName,
			}
			alerts = append(alerts, alert)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	return alerts, nil
}

// PublishAlert publishes an alert to the Pub/Sub topic
func (m *Monitor) PublishAlert(alert *Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	result := m.alertTopic.Publish(m.ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"severity":  alert.Severity,
			"projectId": alert.ProjectID,
		},
	})

	// Wait for the publish result
	_, err = result.Get(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to publish alert: %w", err)
	}

	return nil
}

// Close closes the monitor and releases resources
func (m *Monitor) Close() error {
	m.alertTopic.Stop()
	return m.pubsubClient.Close()
} 