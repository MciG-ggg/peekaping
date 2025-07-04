package providers

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "peekaping/src/modules/heartbeat"
    "peekaping/src/modules/monitor"

    liquid "github.com/osteele/liquid"
    "go.uber.org/zap"
)

// TwilioConfig holds the configuration for Twilio notification channel
// The JSON field names follow the pattern used by the existing UI and API.
//
// Validation rules:
//   - AccountSID must be set (required)
//   - AuthToken must be set (required)
//   - FromNumber and ToNumber must be valid phone numbers; we only ensure they are non-empty here
// The ApiKey and UseTemplate/Template fields are optional.
//
//nolint:revive // field naming kept for backward compatibility with existing payloads
// Use underscore names to match JSON produced by the web application.
type TwilioConfig struct {
    AccountSID  string `json:"twilio_account_sid" validate:"required"`
    ApiKey      string `json:"twilio_api_key"`
    AuthToken   string `json:"twilio_auth_token" validate:"required"`
    FromNumber  string `json:"twilio_from_number" validate:"required"`
    ToNumber    string `json:"twilio_to_number" validate:"required"`
    UseTemplate bool   `json:"use_template"`
    Template    string `json:"template"`
}

// TwilioSender implements NotificationChannelProvider for Twilio SMS.
type TwilioSender struct {
    logger *zap.SugaredLogger
}

// NewTwilioSender returns a new TwilioSender.
func NewTwilioSender(logger *zap.SugaredLogger) *TwilioSender {
    return &TwilioSender{logger: logger}
}

// Unmarshal parses raw JSON config into TwilioConfig.
func (s *TwilioSender) Unmarshal(configJSON string) (any, error) {
    return GenericUnmarshal[TwilioConfig](configJSON)
}

// Validate checks that a configuration JSON string contains the required fields.
func (s *TwilioSender) Validate(configJSON string) error {
    cfg, err := s.Unmarshal(configJSON)
    if err != nil {
        return err
    }
    return GenericValidator(cfg.(*TwilioConfig))
}

// Send constructs and sends the SMS via Twilio REST API.
func (s *TwilioSender) Send(
    ctx context.Context,
    configJSON string,
    message string,
    monitor *monitor.Model,
    heartbeat *heartbeat.Model,
) error {
    cfgAny, err := s.Unmarshal(configJSON)
    if err != nil {
        return err
    }
    cfg := cfgAny.(*TwilioConfig)

    // Render message from template if requested.
    bodyText := message
    if cfg.UseTemplate && cfg.Template != "" {
        engine := liquid.NewEngine()

        bindings := PrepareTemplateBindings(monitor, heartbeat, message)

        // Optional debug output
        if s.logger != nil {
            jsonDebug, _ := json.MarshalIndent(bindings, "", "  ")
            s.logger.Debugf("Template bindings: %s", string(jsonDebug))
        }

        if rendered, err := engine.ParseAndRenderString(cfg.Template, bindings); err == nil {
            bodyText = rendered
        } else {
            return fmt.Errorf("failed to render template: %w", err)
        }
    }

    // Twilio authentication: use API Key if provided, otherwise Account SID.
    authUser := cfg.ApiKey
    if authUser == "" {
        authUser = cfg.AccountSID
    }

    // Prepare form data.
    form := url.Values{}
    form.Set("To", cfg.ToNumber)
    form.Set("From", cfg.FromNumber)
    form.Set("Body", bodyText)

    endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", cfg.AccountSID)

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
    if err != nil {
        return fmt.Errorf("failed to create Twilio request: %w", err)
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.SetBasicAuth(authUser, cfg.AuthToken)
    req.Header.Set("User-Agent", "Peekaping-Twilio/1.0")

    if s.logger != nil {
        s.logger.Debugf("Sending Twilio SMS: %s -> %s", cfg.FromNumber, cfg.ToNumber)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send Twilio message: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("twilio API returned status: %s, body: %s", resp.Status, string(bodyBytes))
    }

    if s.logger != nil {
        s.logger.Infof("Twilio message sent successfully to %s", cfg.ToNumber)
    }

    return nil
}