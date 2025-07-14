package executor

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/shared"
	"strings"
	"time"

	"go.uber.org/zap"
)

type TLSInfo struct {
	CertInfo *CertificateInfo `json:"cert_info,omitempty"`
}

type CertificateInfo struct {
	Subject          string    `json:"subject"`
	Issuer           string    `json:"issuer"`
	NotBefore        time.Time `json:"not_before"`
	NotAfter         time.Time `json:"not_after"`
	DaysRemaining    int       `json:"days_remaining"`
	CertType         string    `json:"cert_type"`
	Fingerprint256   string    `json:"fingerprint_256"`
	IssuerCertificate *CertificateInfo `json:"issuer_certificate,omitempty"`
}

type Result struct {
	Status    heartbeat.MonitorStatus
	Message   string
	StartTime time.Time
	EndTime   time.Time
	TLSInfo   *TLSInfo `json:"tls_info,omitempty"`
}

type Monitor = shared.Monitor
type Proxy = shared.Proxy

// duplicate of monitor.Model to avoid circular dependency
type ExecutorMonitorParams struct {
	ID             string
	Type           string
	Name           string
	Interval       int
	Timeout        int
	MaxRetries     int
	RetryInterval  int
	ResendInterval int
	Active         bool
	Status         heartbeat.MonitorStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Config         string
}

// Executor defines the interface that all health check executors must implement
type Executor interface {
	Execute(ctx context.Context, params *Monitor, proxyModel *Proxy) *Result
	Validate(configJSON string) error
	Unmarshal(configJSON string) (any, error)
}

type ExecutorRegistry struct {
	logger   *zap.SugaredLogger
	registry map[string]Executor
}

func NewExecutorRegistry(logger *zap.SugaredLogger, heartbeatService heartbeat.Service) *ExecutorRegistry {
	registry := make(map[string]Executor)

	registry["http"] = NewHTTPExecutor(logger)
	registry["push"] = NewPushExecutor(logger, heartbeatService)
	registry["tcp"] = NewTCPExecutor(logger)
	registry["ping"] = NewPingExecutor(logger)
	registry["dns"] = NewDNSExecutor(logger)
	registry["docker"] = NewDockerExecutor(logger)
	registry["grpc-keyword"] = NewGRPCExecutor(logger)
	registry["snmp"] = NewSnmpExecutor(logger)

	return &ExecutorRegistry{
		registry: registry,
		logger:   logger,
	}
}

// func (f *ExecutorRegistry) RegisterExecutor(name string, executor Executor) {
// 	f.registry[name] = executor
// }

func (f *ExecutorRegistry) GetExecutor(name string) (Executor, bool) {
	e, ok := f.registry[name]
	return e, ok
}

func (er *ExecutorRegistry) ValidateConfig(monitorType string, configJSON string) error {
	executor, ok := er.GetExecutor(monitorType)
	if !ok {
		err := fmt.Errorf("executor not found for monitor type: %s", monitorType)
		return err
	}

	err := executor.Validate(configJSON)
	if err != nil {
		er.logger.Errorf("failed to validate config: %s", err.Error())
		return err
	}

	return nil
}

// ExtractCertificateInfo extracts certificate information from an x509 certificate
func ExtractCertificateInfo(cert *x509.Certificate) *CertificateInfo {
	if cert == nil {
		return nil
	}

	// Calculate days remaining
	now := time.Now()
	daysRemaining := int(cert.NotAfter.Sub(now).Hours() / 24)

	// Calculate SHA256 fingerprint
	hash := sha256.Sum256(cert.Raw)
	fingerprint := hex.EncodeToString(hash[:])

	// Determine certificate type
	certType := "unknown"
	if cert.IsCA {
		certType = "ca"
	} else {
		certType = "server"
	}

	// Get subject CN
	subject := ""
	if len(cert.Subject.CommonName) > 0 {
		subject = cert.Subject.CommonName
	} else if len(cert.Subject.Organization) > 0 {
		subject = strings.Join(cert.Subject.Organization, ", ")
	}

	// Get issuer CN
	issuer := ""
	if len(cert.Issuer.CommonName) > 0 {
		issuer = cert.Issuer.CommonName
	} else if len(cert.Issuer.Organization) > 0 {
		issuer = strings.Join(cert.Issuer.Organization, ", ")
	}

	return &CertificateInfo{
		Subject:          subject,
		Issuer:           issuer,
		NotBefore:        cert.NotBefore,
		NotAfter:         cert.NotAfter,
		DaysRemaining:    daysRemaining,
		CertType:         certType,
		Fingerprint256:   fingerprint,
		IssuerCertificate: nil, // Will be populated separately if needed
	}
}
