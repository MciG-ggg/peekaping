package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"peekaping/src/modules/shared"
	"peekaping/src/utils"
	"peekaping/src/version"
	"strconv"
	"strings"
	"time"

	"crypto/tls"
	"crypto/x509"

	"github.com/Azure/go-ntlmssp"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

func HTTPConfigStructLevelValidation(sl validator.StructLevel) {
	cfg := sl.Current().Interface().(HTTPConfig)

	switch cfg.Encoding {
	case "json":
		if cfg.Body != "" {
			var js json.RawMessage
			if err := json.Unmarshal([]byte(cfg.Body), &js); err != nil {
				sl.ReportError(cfg.Body, "Body", "body", "json", "")
			}
		}
	case "form":
		if cfg.Body != "" {
			_, err := url.ParseQuery(cfg.Body)
			if err != nil {
				sl.ReportError(cfg.Body, "Body", "body", "form", "")
			}
		}
	case "xml":
		if cfg.Body != "" {
			if err := xml.Unmarshal([]byte(cfg.Body), new(interface{})); err != nil {
				sl.ReportError(cfg.Body, "Body", "body", "xml", "")
			}
		}
	case "text":
		// No validation needed
	}

	// Authentication validation
	switch cfg.AuthMethod {
	case "none":
		// No extra fields required
	case "basic":
		if cfg.BasicAuthUser == "" {
			sl.ReportError(cfg.BasicAuthUser, "BasicAuthUser", "basic_auth_user", "required_with_auth_basic", "")
		}
		if cfg.BasicAuthPass == "" {
			sl.ReportError(cfg.BasicAuthPass, "BasicAuthPass", "basic_auth_pass", "required_with_auth_basic", "")
		}
	case "ntlm":
		if cfg.BasicAuthUser == "" {
			sl.ReportError(cfg.BasicAuthUser, "BasicAuthUser", "basic_auth_user", "required_with_auth_ntlm", "")
		}
		if cfg.BasicAuthPass == "" {
			sl.ReportError(cfg.BasicAuthPass, "BasicAuthPass", "basic_auth_pass", "required_with_auth_ntlm", "")
		}
		if cfg.AuthDomain == "" {
			sl.ReportError(cfg.AuthDomain, "AuthDomain", "authDomain", "required_with_auth_ntlm", "")
		}
		if cfg.AuthWorkstation == "" {
			sl.ReportError(cfg.AuthWorkstation, "AuthWorkstation", "authWorkstation", "required_with_auth_ntlm", "")
		}
	case "oauth2-cc":
		if cfg.OauthAuthMethod != "client_secret_basic" && cfg.OauthAuthMethod != "client_secret_post" {
			sl.ReportError(cfg.OauthAuthMethod, "OauthAuthMethod", "oauth_auth_method", "oneof=client_secret_basic client_secret_post", "")
		}
		if cfg.OauthTokenUrl == "" {
			sl.ReportError(cfg.OauthTokenUrl, "OauthTokenUrl", "oauth_token_url", "required_with_auth_oauth2cc", "")
		} else {
			_, err := url.ParseRequestURI(cfg.OauthTokenUrl)
			if err != nil {
				sl.ReportError(cfg.OauthTokenUrl, "OauthTokenUrl", "oauth_token_url", "url", "")
			}
		}
		if cfg.OauthClientId == "" {
			sl.ReportError(cfg.OauthClientId, "OauthClientId", "oauth_client_id", "required_with_auth_oauth2cc", "")
		}
		if cfg.OauthClientSecret == "" {
			sl.ReportError(cfg.OauthClientSecret, "OauthClientSecret", "oauth_client_secret", "required_with_auth_oauth2cc", "")
		}
		// OauthScopes is optional
	case "mtls":
		if cfg.TlsCert == "" {
			sl.ReportError(cfg.TlsCert, "TlsCert", "tlsCert", "required_with_auth_mtls", "")
		}
		if cfg.TlsKey == "" {
			sl.ReportError(cfg.TlsKey, "TlsKey", "tlsKey", "required_with_auth_mtls", "")
		}
		if cfg.TlsCa == "" {
			sl.ReportError(cfg.TlsCa, "TlsCa", "tlsCa", "required_with_auth_mtls", "")
		}
	}
}

type HTTPConfig struct {
	Url string `json:"url" validate:"required,url"`

	Method              string   `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH HEAD OPTIONS"`
	Headers             string   `json:"headers" validate:"omitempty,json"`
	Encoding            string   `json:"encoding" validate:"required,oneof=json form xml text"`
	Body                string   `json:"body" validate:"omitempty"`
	AcceptedStatusCodes []string `json:"accepted_statuscodes" validate:"required,dive,oneof=2XX 3XX 4XX 5XX"`
	MaxRedirects        int      `json:"max_redirects" validate:"omitempty,min=0"`
	IgnoreTlsErrors     bool     `json:"ignore_tls_errors"`

	// Authentication fields
	AuthMethod        string `json:"authMethod" validate:"required,oneof=none basic oauth2-cc ntlm mtls"`
	BasicAuthUser     string `json:"basic_auth_user,omitempty"`
	BasicAuthPass     string `json:"basic_auth_pass,omitempty"`
	AuthDomain        string `json:"authDomain,omitempty"`
	AuthWorkstation   string `json:"authWorkstation,omitempty"`
	OauthAuthMethod   string `json:"oauth_auth_method,omitempty"`
	OauthTokenUrl     string `json:"oauth_token_url,omitempty"`
	OauthClientId     string `json:"oauth_client_id,omitempty"`
	OauthClientSecret string `json:"oauth_client_secret,omitempty"`
	OauthScopes       string `json:"oauth_scopes,omitempty"`
	TlsCert           string `json:"tlsCert,omitempty"`
	TlsKey            string `json:"tlsKey,omitempty"`
	TlsCa             string `json:"tlsCa,omitempty"`

	// Response validation fields
	ResponseValidation string `json:"response_validation" validate:"omitempty,oneof=none keyword json_query"`
	Keyword           string `json:"keyword,omitempty"`
	InvertKeyword     bool   `json:"invert_keyword,omitempty"`
	JsonQuery         string `json:"json_query,omitempty"`
	JsonQueryCondition string `json:"json_query_condition" validate:"omitempty,oneof=== != > < >= <="`
	JsonQueryExpectedValue string `json:"json_query_expected_value,omitempty"`
}

type HTTPExecutor struct {
	client *http.Client
	logger *zap.SugaredLogger
}

func NewHTTPExecutor(logger *zap.SugaredLogger) *HTTPExecutor {
	utils.Validate.RegisterStructValidation(HTTPConfigStructLevelValidation, HTTPConfig{})

	return &HTTPExecutor{
		client: &http.Client{},
		logger: logger,
	}
}

func (s *HTTPExecutor) Unmarshal(configJSON string) (any, error) {
	return GenericUnmarshal[HTTPConfig](configJSON)
}

func (s *HTTPExecutor) Validate(configJSON string) error {
	cfg, err := s.Unmarshal(configJSON)
	if err != nil {
		return err
	}
	return GenericValidator(cfg.(*HTTPConfig))
}

// Helper to check if status code matches accepted patterns
func isStatusAccepted(statusCode int, accepted []string) bool {
	for _, pattern := range accepted {
		switch pattern {
		case "2XX":
			if statusCode >= 200 && statusCode < 300 {
				return true
			}
		case "3XX":
			if statusCode >= 300 && statusCode < 400 {
				return true
			}
		case "4XX":
			if statusCode >= 400 && statusCode < 500 {
				return true
			}
		case "5XX":
			if statusCode >= 500 && statusCode < 600 {
				return true
			}
		}
	}
	return false
}

func buildProxyTransport(base *http.Transport, proxyModel *Proxy) http.RoundTripper {
	if proxyModel == nil {
		return base
	}

	// Set default protocol if not specified
	protocol := proxyModel.Protocol
	if protocol == "" {
		protocol = "http"
	}

	switch protocol {
	case "http", "https":
		proxyURL := &url.URL{
			Scheme: protocol,
			Host:   fmt.Sprintf("%s:%d", proxyModel.Host, proxyModel.Port),
		}
		if proxyModel.Auth && proxyModel.Username != "" && proxyModel.Password != "" {
			proxyURL.User = url.UserPassword(proxyModel.Username, proxyModel.Password)
		}
		base.Proxy = http.ProxyURL(proxyURL)
		return base
	case "socks", "socks5", "socks5h", "socks4":
		var auth *proxy.Auth
		if proxyModel.Auth && proxyModel.Username != "" && proxyModel.Password != "" {
			auth = &proxy.Auth{
				User:     proxyModel.Username,
				Password: proxyModel.Password,
			}
		}
		address := fmt.Sprintf("%s:%d", proxyModel.Host, proxyModel.Port)
		dialer, err := proxy.SOCKS5("tcp", address, auth, proxy.Direct)
		if err != nil {
			// fallback to default transport if dialer fails
			return base
		}
		base.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			return dialer.Dial(network, addr)
		}
		base.Proxy = nil // No HTTP proxy
		return base
	default:
		return base
	}
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "peekaping/"+version.Version)
	req.Header.Set("Accept", "*/*")
}

func (h *HTTPExecutor) Execute(ctx context.Context, m *Monitor, proxyModel *Proxy) *Result {
	cfgAny, err := h.Unmarshal(m.Config)
	if err != nil {
		return DownResult(err, time.Now().UTC(), time.Now().UTC())
	}
	cfg := cfgAny.(*HTTPConfig)

	h.logger.Debugf("execute http cfg: %+v", cfg)

	var bodyReader io.Reader
	if cfg.Body != "" {
		bodyReader = bytes.NewReader([]byte(cfg.Body))
	}

	req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.Url, bodyReader)
	if err != nil {
		return DownResult(err, time.Now().UTC(), time.Now().UTC())
	}
	setDefaultHeaders(req)

	if cfg.Headers != "" {
		headersMap := make(map[string]string)
		err := json.Unmarshal([]byte(cfg.Headers), &headersMap)
		if err != nil {
			return DownResult(fmt.Errorf("invalid headers json: %w", err), time.Now().UTC(), time.Now().UTC())
		}
		for k, v := range headersMap {
			req.Header.Set(k, v)
		}
	}

	// Determine effective max redirects value
	effectiveMaxRedirects := cfg.MaxRedirects

	checkRedirect := func(req *http.Request, via []*http.Request) error {
		h.logger.Debugf("checkRedirect: %d redirects followed, max allowed: %d", len(via), effectiveMaxRedirects)
		if effectiveMaxRedirects == 0 {
			return fmt.Errorf("redirects disabled: max_redirects set to 0")
		}
		if len(via) > effectiveMaxRedirects {
			return fmt.Errorf("too many redirects: followed %d redirects, maximum allowed is %d", len(via), effectiveMaxRedirects)
		}
		return nil
	}

	switch cfg.Encoding {
	case "json":
		req.Header.Set("Content-Type", "application/json")
	case "form":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "xml":
		req.Header.Set("Content-Type", "application/xml")
	case "text":
		req.Header.Set("Content-Type", "text/plain")
	}

	// --- PROXY LOGIC ---

	// Default transport with proxy if needed
	baseTransport := &http.Transport{}

	// Configure TLS settings if needed
	if cfg.IgnoreTlsErrors {
		if baseTransport.TLSClientConfig == nil {
			baseTransport.TLSClientConfig = &tls.Config{}
		}
		baseTransport.TLSClientConfig.InsecureSkipVerify = true
	}

	transport := buildProxyTransport(baseTransport, proxyModel)

	// Set timeout from monitor configuration
	timeout := time.Duration(m.Timeout) * time.Second

	// --- AUTHENTICATION LOGIC ---
	switch cfg.AuthMethod {
	case "basic":
		req.SetBasicAuth(cfg.BasicAuthUser, cfg.BasicAuthPass)
	case "ntlm":
		// NTLM authentication using github.com/Azure/go-ntlmssp
		ntlmTransport := ntlmssp.Negotiator{
			RoundTripper: transport,
		}
		h.client = &http.Client{
			Transport:     &ntlmTransport,
			Timeout:       time.Duration(m.Timeout) * time.Second,
			CheckRedirect: checkRedirect,
		}

		if cfg.AuthDomain != "" {
			req.SetBasicAuth(cfg.AuthDomain+"\\"+cfg.BasicAuthUser, cfg.BasicAuthPass)
		} else {
			req.SetBasicAuth(cfg.BasicAuthUser, cfg.BasicAuthPass)
		}
	case "oauth2-cc":
		form := url.Values{}
		form.Set("grant_type", "client_credentials")
		if cfg.OauthScopes != "" {
			form.Set("scope", cfg.OauthScopes)
		}
		form.Set("client_id", cfg.OauthClientId)
		form.Set("client_secret", cfg.OauthClientSecret)

		tokenReq, err := http.NewRequestWithContext(ctx, "POST", cfg.OauthTokenUrl, strings.NewReader(form.Encode()))
		if err != nil {
			return DownResult(fmt.Errorf("failed to create oauth2 token request: %w", err), time.Now().UTC(), time.Now().UTC())
		}
		setDefaultHeaders(tokenReq)

		tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if cfg.OauthAuthMethod == "client_secret_basic" {
			basic := base64.StdEncoding.EncodeToString([]byte(cfg.OauthClientId + ":" + cfg.OauthClientSecret))
			tokenReq.Header.Set("Authorization", "Basic "+basic)
		}

		tokenResp, err := http.DefaultClient.Do(tokenReq)
		if err != nil {
			return DownResult(fmt.Errorf("failed to get oauth2 token: %w", err), time.Now().UTC(), time.Now().UTC())
		}
		defer tokenResp.Body.Close()
		if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
			return DownResult(fmt.Errorf("oauth2 token endpoint returned status: %d", tokenResp.StatusCode), time.Now().UTC(), time.Now().UTC())
		}
		var tokenData struct {
			AccessToken string `json:"access_token"`
		}
		err = json.NewDecoder(tokenResp.Body).Decode(&tokenData)
		if err != nil || tokenData.AccessToken == "" {
			return DownResult(fmt.Errorf("failed to parse oauth2 token response: %w", err), time.Now().UTC(), time.Now().UTC())
		}
		req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	case "mtls":
		cert, err := tls.X509KeyPair([]byte(cfg.TlsCert), []byte(cfg.TlsKey))
		if err != nil {
			return DownResult(fmt.Errorf("invalid mTLS cert/key: %w", err), time.Now().UTC(), time.Now().UTC())
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM([]byte(cfg.TlsCa)); !ok {
			return DownResult(fmt.Errorf("invalid mTLS CA cert"), time.Now().UTC(), time.Now().UTC())
		}
		mtlsTransport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            caCertPool,
				InsecureSkipVerify: cfg.IgnoreTlsErrors,
			},
		}
		mtlsTransportWithProxy := buildProxyTransport(mtlsTransport, proxyModel)
		h.client = &http.Client{
			Transport:     mtlsTransportWithProxy,
			Timeout:       time.Duration(m.Timeout) * time.Second,
			CheckRedirect: checkRedirect,
		}
	}

	if cfg.AuthMethod != "mtls" && cfg.AuthMethod != "ntlm" {
		h.client = &http.Client{
			Timeout:       timeout,
			CheckRedirect: checkRedirect,
			Transport:     transport,
		}
	}

	// Set user agent and accept headers

	startTime := time.Now().UTC()
	resp, err := h.client.Do(req)
	endTime := time.Now().UTC()

	if err != nil {
		h.logger.Infof("HTTP request failed: %s, %s", m.Name, err.Error())
		return DownResult(err, startTime, endTime)
	}
	defer resp.Body.Close()

	h.logger.Infof("HTTP response status: %s, %d", m.Name, resp.StatusCode)

	if !isStatusAccepted(resp.StatusCode, cfg.AcceptedStatusCodes) {
		return &Result{
			Status:    shared.MonitorStatusDown,
			Message:   fmt.Sprintf("HTTP request failed with status: %d", resp.StatusCode),
			StartTime: startTime,
			EndTime:   endTime,
		}
	}

	// Read response body for content validation
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Warnf("Failed to read response body: %v", err)
		responseBody = []byte{}
	}

	// Validate response content based on configured validation mode
	if cfg.ResponseValidation != "" && cfg.ResponseValidation != "none" {
		validationResult := h.validateResponseContent(cfg, responseBody)
		if validationResult.Status == shared.MonitorStatusDown {
			return &Result{
				Status:    shared.MonitorStatusDown,
				Message:   validationResult.Message,
				StartTime: startTime,
				EndTime:   endTime,
			}
		}
	}

	return &Result{
		Status:    shared.MonitorStatusUp,
		Message:   fmt.Sprintf("%d - %s", resp.StatusCode, resp.Status),
		StartTime: startTime,
		EndTime:   endTime,
	}
}

func (h *HTTPExecutor) validateResponseContent(cfg *HTTPConfig, responseBody []byte) *Result {
	responseText := string(responseBody)
	
	switch cfg.ResponseValidation {
	case "keyword":
		return h.validateKeyword(cfg, responseText)
	case "json_query":
		return h.validateJsonQuery(cfg, responseBody)
	default:
		return &Result{Status: shared.MonitorStatusUp, Message: "No validation configured"}
	}
}

func (h *HTTPExecutor) validateKeyword(cfg *HTTPConfig, responseText string) *Result {
	if cfg.Keyword == "" {
		return &Result{Status: shared.MonitorStatusUp, Message: "No keyword specified"}
	}

	containsKeyword := strings.Contains(responseText, cfg.Keyword)
	
	if cfg.InvertKeyword {
		// Invert keyword: should NOT contain the keyword
		if containsKeyword {
			return &Result{
				Status:  shared.MonitorStatusDown,
				Message: fmt.Sprintf("Keyword '%s' found in response (inverted check)", cfg.Keyword),
			}
		}
		return &Result{
			Status:  shared.MonitorStatusUp,
			Message: fmt.Sprintf("Keyword '%s' not found in response (inverted check)", cfg.Keyword),
		}
	} else {
		// Normal keyword: should contain the keyword
		if !containsKeyword {
			return &Result{
				Status:  shared.MonitorStatusDown,
				Message: fmt.Sprintf("Keyword '%s' not found in response", cfg.Keyword),
			}
		}
		return &Result{
			Status:  shared.MonitorStatusUp,
			Message: fmt.Sprintf("Keyword '%s' found in response", cfg.Keyword),
		}
	}
}

func (h *HTTPExecutor) validateJsonQuery(cfg *HTTPConfig, responseBody []byte) *Result {
	if cfg.JsonQuery == "" {
		return &Result{Status: shared.MonitorStatusUp, Message: "No JSON query specified"}
	}

	// Handle raw response query
	if cfg.JsonQuery == "$" {
		responseStr := string(responseBody)
		if cfg.JsonQueryCondition == "" || cfg.JsonQueryExpectedValue == "" {
			return &Result{
				Status:  shared.MonitorStatusUp,
				Message: fmt.Sprintf("Raw response: %s", responseStr),
			}
		}
		return h.compareJsonQueryResult(responseStr, cfg.JsonQueryCondition, cfg.JsonQueryExpectedValue)
	}

	// Parse JSON response
	var jsonData interface{}
	if err := json.Unmarshal(responseBody, &jsonData); err != nil {
		return &Result{
			Status:  shared.MonitorStatusDown,
			Message: fmt.Sprintf("Invalid JSON response: %v", err),
		}
	}

	// Simple JSON path implementation (only supports basic key access)
	result, err := h.extractJsonValue(jsonData, cfg.JsonQuery)
	if err != nil {
		return &Result{
			Status:  shared.MonitorStatusDown,
			Message: fmt.Sprintf("JSON query failed: %v", err),
		}
	}

	// If no condition is specified, just check if the result exists
	if cfg.JsonQueryCondition == "" || cfg.JsonQueryExpectedValue == "" {
		if result == nil {
			return &Result{
				Status:  shared.MonitorStatusDown,
				Message: "JSON query returned null/undefined",
			}
		}
		return &Result{
			Status:  shared.MonitorStatusUp,
			Message: fmt.Sprintf("JSON query returned: %v", result),
		}
	}

	// Compare result with expected value
	return h.compareJsonQueryResult(result, cfg.JsonQueryCondition, cfg.JsonQueryExpectedValue)
}

func (h *HTTPExecutor) compareJsonQueryResult(result interface{}, condition, expectedValue string) *Result {
	// Convert result to string for comparison
	resultStr := fmt.Sprintf("%v", result)
	
	switch condition {
	case "===":
		if resultStr == expectedValue {
			return &Result{Status: shared.MonitorStatusUp, Message: fmt.Sprintf("JSONata result '%s' equals expected value '%s'", resultStr, expectedValue)}
		}
		return &Result{Status: shared.MonitorStatusDown, Message: fmt.Sprintf("JSONata result '%s' does not equal expected value '%s'", resultStr, expectedValue)}
	
	case "!=":
		if resultStr != expectedValue {
			return &Result{Status: shared.MonitorStatusUp, Message: fmt.Sprintf("JSONata result '%s' does not equal expected value '%s'", resultStr, expectedValue)}
		}
		return &Result{Status: shared.MonitorStatusDown, Message: fmt.Sprintf("JSONata result '%s' equals expected value '%s' (should not)", resultStr, expectedValue)}
	
	case ">", "<", ">=", "<=":
		// Try to convert to numbers for numeric comparison
		resultNum, err1 := strconv.ParseFloat(resultStr, 64)
		expectedNum, err2 := strconv.ParseFloat(expectedValue, 64)
		
		if err1 != nil || err2 != nil {
			return &Result{Status: shared.MonitorStatusDown, Message: fmt.Sprintf("Cannot compare non-numeric values: result='%s', expected='%s'", resultStr, expectedValue)}
		}
		
		var comparison bool
		switch condition {
		case ">":
			comparison = resultNum > expectedNum
		case "<":
			comparison = resultNum < expectedNum
		case ">=":
			comparison = resultNum >= expectedNum
		case "<=":
			comparison = resultNum <= expectedNum
		}
		
		if comparison {
			return &Result{Status: shared.MonitorStatusUp, Message: fmt.Sprintf("JSONata result %s %s %s", resultStr, condition, expectedValue)}
		}
		return &Result{Status: shared.MonitorStatusDown, Message: fmt.Sprintf("JSONata result %s is not %s %s", resultStr, condition, expectedValue)}
	
	default:
		return &Result{Status: shared.MonitorStatusDown, Message: fmt.Sprintf("Unknown condition: %s", condition)}
	}
}

// extractJsonValue extracts a value from JSON data using a simple path
// Supports basic dot notation like "user.name" or "data.items[0].id"
func (h *HTTPExecutor) extractJsonValue(data interface{}, path string) (interface{}, error) {
	if path == "" {
		return data, nil
	}

	// Split path by dots
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		// Handle array indexing like "items[0]"
		if strings.Contains(key, "[") && strings.Contains(key, "]") {
			// Extract key and index
			parts := strings.Split(key, "[")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid array syntax: %s", key)
			}
			arrayKey := parts[0]
			indexStr := strings.TrimSuffix(parts[1], "]")
			
			// Get the array
			current = h.getObjectValue(current, arrayKey)
			if current == nil {
				return nil, fmt.Errorf("key not found: %s", arrayKey)
			}
			
			// Convert to slice
			slice, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("not an array: %s", arrayKey)
			}
			
			// Parse index
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", indexStr)
			}
			
			if index < 0 || index >= len(slice) {
				return nil, fmt.Errorf("array index out of bounds: %d", index)
			}
			
			current = slice[index]
		} else {
			// Regular key access
			current = h.getObjectValue(current, key)
			if current == nil {
				return nil, fmt.Errorf("key not found: %s", key)
			}
		}
	}

	return current, nil
}

// getObjectValue gets a value from an object by key
func (h *HTTPExecutor) getObjectValue(data interface{}, key string) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case map[string]interface{}:
		return v[key]
	case map[interface{}]interface{}:
		return v[key]
	default:
		return nil
	}
}
