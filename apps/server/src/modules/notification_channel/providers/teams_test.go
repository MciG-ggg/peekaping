package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"peekaping/src/config"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/shared"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTeamsSender_Unmarshal(t *testing.T) {
	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	tests := []struct {
		name        string
		configJSON  string
		expected    *TeamsConfig
		expectError bool
	}{
		{
			name:       "valid config",
			configJSON: `{"webhook_url":"https://example.com/webhook","use_template":true,"template":"Custom template"}`,
			expected: &TeamsConfig{
				WebhookURL:  "https://example.com/webhook",
				UseTemplate: true,
				Template:    "Custom template",
			},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			configJSON:  `{"webhook_url":}`,
			expectError: true,
		},
		{
			name:       "minimal config",
			configJSON: `{"webhook_url":"https://example.com/webhook"}`,
			expected: &TeamsConfig{
				WebhookURL:  "https://example.com/webhook",
				UseTemplate: false,
				Template:    "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sender.Unmarshal(tt.configJSON)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			cfg := result.(*TeamsConfig)
			assert.Equal(t, tt.expected.WebhookURL, cfg.WebhookURL)
			assert.Equal(t, tt.expected.UseTemplate, cfg.UseTemplate)
			assert.Equal(t, tt.expected.Template, cfg.Template)
		})
	}
}

func TestTeamsSender_Validate(t *testing.T) {
	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	tests := []struct {
		name        string
		configJSON  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config",
			configJSON:  `{"webhook_url":"https://example.com/webhook"}`,
			expectError: false,
		},
		{
			name:        "missing webhook_url",
			configJSON:  `{"use_template":true}`,
			expectError: true,
			errorMsg:    "WebhookURL",
		},
		{
			name:        "invalid webhook_url format",
			configJSON:  `{"webhook_url":"not-a-url"}`,
			expectError: true,
			errorMsg:    "WebhookURL",
		},
		{
			name:        "empty webhook_url",
			configJSON:  `{"webhook_url":""}`,
			expectError: true,
			errorMsg:    "WebhookURL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sender.Validate(tt.configJSON)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTeamsSender_StatusMessageFactory(t *testing.T) {
	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	tests := []struct {
		name             string
		status           heartbeat.MonitorStatus
		monitorName      string
		withStatusSymbol bool
		expected         string
	}{
		{
			name:             "down status with symbol",
			status:           shared.MonitorStatusDown,
			monitorName:      "Test Monitor",
			withStatusSymbol: true,
			expected:         "🔴 [Test Monitor] went down",
		},
		{
			name:             "down status without symbol",
			status:           shared.MonitorStatusDown,
			monitorName:      "Test Monitor",
			withStatusSymbol: false,
			expected:         "[Test Monitor] went down",
		},
		{
			name:             "up status with symbol",
			status:           shared.MonitorStatusUp,
			monitorName:      "Test Monitor",
			withStatusSymbol: true,
			expected:         "✅ [Test Monitor] is back online",
		},
		{
			name:             "up status without symbol",
			status:           shared.MonitorStatusUp,
			monitorName:      "Test Monitor",
			withStatusSymbol: false,
			expected:         "[Test Monitor] is back online",
		},
		{
			name:             "pending status",
			status:           shared.MonitorStatusPending,
			monitorName:      "Test Monitor",
			withStatusSymbol: true,
			expected:         "Notification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sender.statusMessageFactory(tt.status, tt.monitorName, tt.withStatusSymbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTeamsSender_GetStyle(t *testing.T) {
	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	tests := []struct {
		name     string
		status   heartbeat.MonitorStatus
		expected string
	}{
		{
			name:     "down status",
			status:   shared.MonitorStatusDown,
			expected: "attention",
		},
		{
			name:     "up status",
			status:   shared.MonitorStatusUp,
			expected: "good",
		},
		{
			name:     "pending status",
			status:   shared.MonitorStatusPending,
			expected: "emphasis",
		},
		{
			name:     "maintenance status",
			status:   shared.MonitorStatusMaintenance,
			expected: "emphasis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sender.getStyle(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTeamsSender_Send(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Peekaping-Teams/1.0", r.Header.Get("User-Agent"))

		// Verify payload structure
		var payload map[string]any
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		// Check basic structure
		assert.Equal(t, "message", payload["type"])
		assert.NotEmpty(t, payload["summary"])

		attachments, ok := payload["attachments"].([]any)
		assert.True(t, ok)
		assert.Len(t, attachments, 1)

		attachment := attachments[0].(map[string]any)
		assert.Equal(t, "application/vnd.microsoft.card.adaptive", attachment["contentType"])

		content := attachment["content"].(map[string]any)
		assert.Equal(t, "AdaptiveCard", content["type"])
		assert.Equal(t, "1.5", content["version"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	}))
	defer server.Close()

	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{
		ClientURL: "https://peekaping.example.com",
	})

	configJSON := `{"webhook_url":"` + server.URL + `"}`
	message := "Test notification"

	monitor := &monitor.Model{
		ID:   "test-monitor-id",
		Name: "Test Monitor",
		Type: "http",
	}

	heartbeat := &heartbeat.Model{
		Status: shared.MonitorStatusDown,
		Msg:    "Monitor is down",
		Time:   time.Now(),
	}

	err := sender.Send(context.Background(), configJSON, message, monitor, heartbeat)
	assert.NoError(t, err)
}

func TestTeamsSender_Send_WithTemplate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		// Verify template was processed - check the facts in the adaptive card for rendered template
		attachments := payload["attachments"].([]any)
		attachment := attachments[0].(map[string]any)
		content := attachment["content"].(map[string]any)
		bodyItems := content["body"].([]any)

		// Find the FactSet and check for template content
		var facts []any
		for _, item := range bodyItems {
			itemMap := item.(map[string]any)
			if itemMap["type"] == "FactSet" {
				facts = itemMap["facts"].([]any)
				break
			}
		}

		// Check that facts contain template-rendered content
		found := false
		for _, fact := range facts {
			factMap := fact.(map[string]any)
			if factMap["title"] == "Description" {
				value := factMap["value"].(string)
				assert.Contains(t, value, "Test Monitor")
				assert.Contains(t, value, "DOWN")
				found = true
				break
			}
		}
		assert.True(t, found, "Template-rendered message should be in description fact")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	}))
	defer server.Close()

	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	template := "Monitor {{monitor.name}} is {{status}}"
	configJSON := `{"webhook_url":"` + server.URL + `","use_template":true,"template":"` + template + `"}`

	monitor := &monitor.Model{
		ID:   "test-monitor-id",
		Name: "Test Monitor",
		Type: "http",
	}

	heartbeat := &heartbeat.Model{
		Status: shared.MonitorStatusDown,
		Msg:    "Monitor is down",
		Time:   time.Now(),
	}

	err := sender.Send(context.Background(), configJSON, "Original message", monitor, heartbeat)
	assert.NoError(t, err)
}

func TestTeamsSender_Send_GeneralNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		// Verify it's a general notification (no specific monitor info)
		summary := payload["summary"].(string)
		assert.Contains(t, summary, "Notification")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	}))
	defer server.Close()

	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	configJSON := `{"webhook_url":"` + server.URL + `"}`
	message := "General test notification"

	err := sender.Send(context.Background(), configJSON, message, nil, nil)
	assert.NoError(t, err)
}

func TestTeamsSender_Send_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	configJSON := `{"webhook_url":"` + server.URL + `"}`
	message := "Test notification"

	monitor := &monitor.Model{
		ID:   "test-monitor-id",
		Name: "Test Monitor",
		Type: "http",
	}

	heartbeat := &heartbeat.Model{
		Status: shared.MonitorStatusDown,
		Msg:    "Monitor is down",
		Time:   time.Now(),
	}

	err := sender.Send(context.Background(), configJSON, message, monitor, heartbeat)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Teams webhook returned status")
}

func TestTeamsSender_Send_InvalidConfig(t *testing.T) {
	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	configJSON := `{"invalid": "config"}`
	message := "Test notification"

	err := sender.Send(context.Background(), configJSON, message, nil, nil)
	assert.Error(t, err)
	// The error should be about sending the webhook because validation isn't done in Send method
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

func TestTeamsSender_NotificationPayloadFactory(t *testing.T) {
	sender := NewTeamsSender(zap.NewNop().Sugar(), &config.Config{})

	heartbeat := &heartbeat.Model{
		Status: shared.MonitorStatusDown,
		Msg:    "Test message",
		Time:   time.Now(),
	}

	monitorName := "Test Monitor"
	monitorURL := "https://example.com"
	dashboardURL := "https://peekaping.example.com/monitors/123"
	message := "Custom message"

	payload := sender.notificationPayloadFactory(heartbeat, monitorName, monitorURL, dashboardURL, message)

	// Verify basic structure
	assert.Equal(t, "message", payload["type"])
	assert.NotEmpty(t, payload["summary"])

	attachments := payload["attachments"].([]map[string]any)
	assert.Len(t, attachments, 1)

	attachment := attachments[0]
	assert.Equal(t, "application/vnd.microsoft.card.adaptive", attachment["contentType"])

	content := attachment["content"].(map[string]any)
	assert.Equal(t, "AdaptiveCard", content["type"])
	assert.Equal(t, "1.5", content["version"])

	bodyItems := content["body"].([]map[string]any)
	assert.GreaterOrEqual(t, len(bodyItems), 2) // Should have at least container and fact set

	// Verify facts are included
	factSetItem := bodyItems[1]
	assert.Equal(t, "FactSet", factSetItem["type"])
	facts := factSetItem["facts"].([]map[string]any)
	assert.Greater(t, len(facts), 0)

	// Should have actions since we provided URLs
	assert.Len(t, bodyItems, 3) // Container, FactSet, ActionSet
	actionSetItem := bodyItems[2]
	assert.Equal(t, "ActionSet", actionSetItem["type"])
	actions := actionSetItem["actions"].([]map[string]any)
	assert.Len(t, actions, 2) // Dashboard and monitor URL actions
}
