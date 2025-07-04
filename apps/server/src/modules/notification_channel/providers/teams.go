package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"peekaping/src/config"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/shared"
	"strings"
	"time"

	liquid "github.com/osteele/liquid"
	"go.uber.org/zap"
)

type TeamsConfig struct {
	WebhookURL  string `json:"webhook_url" validate:"required,url"`
	UseTemplate bool   `json:"use_template"`
	Template    string `json:"template"`
}

type TeamsSender struct {
	logger *zap.SugaredLogger
	config *config.Config
	client *http.Client
}

// NewTeamsSender creates a TeamsSender
func NewTeamsSender(logger *zap.SugaredLogger, config *config.Config) *TeamsSender {
	return &TeamsSender{
		logger: logger,
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (t *TeamsSender) Unmarshal(configJSON string) (any, error) {
	return GenericUnmarshal[TeamsConfig](configJSON)
}

func (t *TeamsSender) Validate(configJSON string) error {
	cfg, err := t.Unmarshal(configJSON)
	if err != nil {
		return err
	}
	return GenericValidator(cfg.(*TeamsConfig))
}

// statusMessageFactory generates the message to send based on status
func (t *TeamsSender) statusMessageFactory(status heartbeat.MonitorStatus, monitorName string, withStatusSymbol bool) string {
	if status == shared.MonitorStatusDown {
		prefix := ""
		if withStatusSymbol {
			prefix = "🔴 "
		}
		return fmt.Sprintf("%s[%s] went down", prefix, monitorName)
	} else if status == shared.MonitorStatusUp {
		prefix := ""
		if withStatusSymbol {
			prefix = "✅ "
		}
		return fmt.Sprintf("%s[%s] is back online", prefix, monitorName)
	}
	return "Notification"
}

// getStyle selects the style to use based on status
func (t *TeamsSender) getStyle(status heartbeat.MonitorStatus) string {
	switch status {
	case shared.MonitorStatusDown:
		return "attention"
	case shared.MonitorStatusUp:
		return "good"
	default:
		return "emphasis"
	}
}

// extractAddress extracts the URL from monitor
func (t *TeamsSender) extractAddress(monitor *monitor.Model) string {
	if monitor == nil {
		return ""
	}

	// Try to extract URL from monitor config JSON
	// This would need to be enhanced based on the actual config structure
	// For now, return empty string as the config structure varies by monitor type
	return ""
}

// notificationPayloadFactory generates payload for notification
func (t *TeamsSender) notificationPayloadFactory(heartbeat *heartbeat.Model, monitorName, monitorURL, dashboardURL, message string) map[string]any {
	status := shared.MonitorStatusDown
	if heartbeat != nil {
		status = heartbeat.Status
	}

	facts := []map[string]any{}
	actions := []map[string]any{}

	// Add dashboard URL action if available
	if dashboardURL != "" {
		actions = append(actions, map[string]any{
			"type":  "Action.OpenUrl",
			"title": "Visit Peekaping",
			"url":   dashboardURL,
		})
	}

	// Add message fact if available
	if message != "" {
		facts = append(facts, map[string]any{
			"title": "Description",
			"value": message,
		})
	} else if heartbeat != nil && heartbeat.Msg != "" {
		facts = append(facts, map[string]any{
			"title": "Description",
			"value": heartbeat.Msg,
		})
	}

	// Add monitor name fact
	if monitorName != "" {
		facts = append(facts, map[string]any{
			"title": "Monitor",
			"value": monitorName,
		})
	}

	// Add monitor URL fact and action if available
	if monitorURL != "" && monitorURL != "https://" {
		facts = append(facts, map[string]any{
			"title": "URL",
			"value": fmt.Sprintf("[%s](%s)", monitorURL, monitorURL),
		})
		actions = append(actions, map[string]any{
			"type":  "Action.OpenUrl",
			"title": "Visit Monitor URL",
			"url":   monitorURL,
		})
	}

	// Add time fact if available
	if heartbeat != nil {
		facts = append(facts, map[string]any{
			"title": "Time",
			"value": heartbeat.Time.Format("2006-01-02 15:04:05 MST"),
		})
	}

	// Build the adaptive card
	body := []map[string]any{
		{
			"type":                     "Container",
			"verticalContentAlignment": "Center",
			"items": []map[string]any{
				{
					"type":  "ColumnSet",
					"style": t.getStyle(status),
					"columns": []map[string]any{
						{
							"type":                     "Column",
							"width":                    "auto",
							"verticalContentAlignment": "Center",
							"items": []map[string]any{
								{
									"type":    "Image",
									"width":   "32px",
									"style":   "Person",
									"url":     "https://raw.githubusercontent.com/louislam/uptime-kuma/master/public/icon.png",
									"altText": "Peekaping Logo",
								},
							},
						},
						{
							"type":  "Column",
							"width": "stretch",
							"items": []map[string]any{
								{
									"type":   "TextBlock",
									"size":   "Medium",
									"weight": "Bolder",
									"text":   fmt.Sprintf("**%s**", t.statusMessageFactory(status, monitorName, false)),
								},
								{
									"type":     "TextBlock",
									"size":     "Small",
									"weight":   "Default",
									"text":     "Peekaping Alert",
									"isSubtle": true,
									"spacing":  "None",
								},
							},
						},
					},
				},
			},
		},
		{
			"type":      "FactSet",
			"separator": false,
			"facts":     facts,
		},
	}

	// Add actions if available
	if len(actions) > 0 {
		body = append(body, map[string]any{
			"type":    "ActionSet",
			"actions": actions,
		})
	}

	payload := map[string]any{
		"type":    "message",
		"summary": t.statusMessageFactory(status, monitorName, true),
		"attachments": []map[string]any{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"contentUrl":  "",
				"content": map[string]any{
					"type":    "AdaptiveCard",
					"body":    body,
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"version": "1.5",
				},
			},
		},
	}

	return payload
}

// sendNotification sends the notification via HTTP
func (t *TeamsSender) sendNotification(ctx context.Context, webhookURL string, payload map[string]any) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Teams payload: %w", err)
	}

	t.logger.Debugf("Teams payload: %s", string(jsonPayload))

	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Peekaping-Teams/1.0")

	t.logger.Debugf("Sending Teams webhook request: %s", req.URL.String())

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Teams webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams webhook returned status: %s", resp.Status)
	}

	return nil
}

// handleGeneralNotification sends a general notification
func (t *TeamsSender) handleGeneralNotification(ctx context.Context, webhookURL, msg string) error {
	payload := t.notificationPayloadFactory(nil, "", "", "", msg)
	// Override the summary for general notifications
	payload["summary"] = "Notification"

	// Update the description fact to use the message
	attachments := payload["attachments"].([]map[string]any)
	content := attachments[0]["content"].(map[string]any)
	body := content["body"].([]map[string]any)

	// Find and update the FactSet
	for i, item := range body {
		if item["type"] == "FactSet" {
			// Replace all facts with just the message
			body[i]["facts"] = []map[string]any{
				{
					"title": "Message",
					"value": msg,
				},
			}
			break
		}
	}

	return t.sendNotification(ctx, webhookURL, payload)
}

func (t *TeamsSender) Send(
	ctx context.Context,
	configJSON string,
	message string,
	monitor *monitor.Model,
	heartbeat *heartbeat.Model,
) error {
	cfgAny, err := t.Unmarshal(configJSON)
	if err != nil {
		return err
	}
	cfg := cfgAny.(*TeamsConfig)

	t.logger.Infof("Sending Teams message to webhook: %s", cfg.WebhookURL)

	// Prepare template bindings if using template
	if cfg.UseTemplate && cfg.Template != "" {
		bindings := PrepareTemplateBindings(monitor, heartbeat, message)
		engine := liquid.NewEngine()
		if rendered, err := engine.ParseAndRenderString(cfg.Template, bindings); err == nil {
			message = rendered
		} else {
			return fmt.Errorf("failed to render template: %w", err)
		}
	}

	// Handle general notification (no heartbeat)
	if heartbeat == nil {
		return t.handleGeneralNotification(ctx, cfg.WebhookURL, message)
	}

	// Build URLs
	var dashboardURL string
	monitorName := "Unknown Monitor"
	if monitor != nil {
		monitorName = monitor.Name
		if t.config.ClientURL != "" {
			dashboardURL = fmt.Sprintf("%s/monitors/%s", strings.TrimRight(t.config.ClientURL, "/"), monitor.ID)
		}
	}

	monitorURL := t.extractAddress(monitor)

	// Generate payload and send
	payload := t.notificationPayloadFactory(heartbeat, monitorName, monitorURL, dashboardURL, message)
	err = t.sendNotification(ctx, cfg.WebhookURL, payload)
	if err != nil {
		return err
	}

	t.logger.Infof("Teams message sent successfully")
	return nil
}
