package certificate_expiry

import (
	"context"
	"fmt"
	"peekaping/src/modules/healthcheck/executor"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/notification_channel"
	"peekaping/src/modules/notification_sent_history"
	"peekaping/src/modules/setting"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Service interface {
	CheckCertExpiryNotifications(ctx context.Context, result *executor.Result, m *monitor.Model) error
}

type ServiceImpl struct {
	settingService               setting.Service
	notificationService          notification_channel.Service
	notificationSentHistoryService notification_sent_history.Service
	logger                       *zap.SugaredLogger
}

func NewService(
	settingService setting.Service,
	notificationService notification_channel.Service,
	notificationSentHistoryService notification_sent_history.Service,
	logger *zap.SugaredLogger,
) Service {
	return &ServiceImpl{
		settingService:               settingService,
		notificationService:          notificationService,
		notificationSentHistoryService: notificationSentHistoryService,
		logger:                       logger.Named("[certificate-expiry-service]"),
	}
}

func (s *ServiceImpl) CheckCertExpiryNotifications(ctx context.Context, result *executor.Result, m *monitor.Model) error {
	// Only process if TLS info is available and has certificate info
	if result.TLSInfo == nil || result.TLSInfo.CertInfo == nil {
		s.logger.Debugf("No TLS certificate info available for monitor %s", m.Name)
		return nil
	}

	s.logger.Debugf("Checking certificate expiry for monitor %s", m.Name)

	// Get notification threshold days setting
	notifyDays, err := s.getNotificationDays(ctx)
	if err != nil {
		s.logger.Errorf("Failed to get notification days setting: %v", err)
		return err
	}

	// Process certificate chain
	certInfo := result.TLSInfo.CertInfo
	for certInfo != nil {
		// Skip root certificates (you may want to maintain a list of known root CAs)
		subjectCN := certInfo.Subject
		if s.isRootCertificate(certInfo) {
			s.logger.Debugf("Known root cert: %s certificate \"%s\" (%d days valid), skipping", 
				certInfo.CertType, subjectCN, certInfo.DaysRemaining)
			break
		}

		// Check each notification threshold
		for _, targetDays := range notifyDays {
			if certInfo.DaysRemaining > targetDays {
				s.logger.Debugf("No need to send cert notification for %s certificate \"%s\" (%d days valid) on %d deadline",
					certInfo.CertType, subjectCN, certInfo.DaysRemaining, targetDays)
			} else {
				s.logger.Debugf("Calling sendCertNotificationByTargetDays for %d deadline on certificate %s",
					targetDays, subjectCN)
				err := s.sendCertNotificationByTargetDays(ctx, m, certInfo, targetDays)
				if err != nil {
					s.logger.Errorf("Failed to send certificate notification: %v", err)
				}
			}
		}

		// Move to next certificate in chain
		certInfo = certInfo.IssuerCertificate
	}

	return nil
}

func (s *ServiceImpl) getNotificationDays(ctx context.Context) ([]int, error) {
	// Get the setting for TLS expiry notification days
	setting, err := s.settingService.GetByKey(ctx, "TLS_EXPIRY_NOTIFY_DAYS")
	if err != nil {
		return nil, err
	}

	// Set default if setting doesn't exist
	if setting == nil {
		defaultDays := []int{7, 14, 21}
		// Create the setting with default values
		_, err = s.settingService.SetByKey(ctx, "TLS_EXPIRY_NOTIFY_DAYS", &setting.CreateUpdateDto{
			Value: "7,14,21",
			Type:  "string",
		})
		if err != nil {
			s.logger.Errorf("Failed to create default TLS_EXPIRY_NOTIFY_DAYS setting: %v", err)
		}
		return defaultDays, nil
	}

	// Parse the comma-separated string into int slice
	daysStr := strings.Split(setting.Value, ",")
	var days []int
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			s.logger.Warnf("Invalid day value in TLS_EXPIRY_NOTIFY_DAYS: %s", dayStr)
			continue
		}
		days = append(days, day)
	}

	if len(days) == 0 {
		// Fallback to default if no valid days found
		days = []int{7, 14, 21}
	}

	return days, nil
}

func (s *ServiceImpl) sendCertNotificationByTargetDays(ctx context.Context, m *monitor.Model, certInfo *executor.CertificateInfo, targetDays int) error {
	// Check if notification was already sent for this threshold
	sent, err := s.notificationSentHistoryService.FindByTypeMonitorAndDays(ctx, "certificate", m.ID, targetDays)
	if err != nil {
		return fmt.Errorf("failed to check notification history: %w", err)
	}

	if sent != nil {
		s.logger.Debugf("Certificate notification already sent for monitor %s, days %d", m.ID, targetDays)
		return nil
	}

	// Prepare notification message
	message := fmt.Sprintf("[%s][%s] %s certificate %s will expire in %d days",
		m.Name, m.Config, certInfo.CertType, certInfo.Subject, certInfo.DaysRemaining)

	s.logger.Debugf("Sending certificate notification: %s", message)

	// TODO: Send notification through notification channels
	// This would require extending the notification system to handle certificate notifications
	// For now, we'll just log and record that we "sent" the notification

	// Record that notification was sent
	_, err = s.notificationSentHistoryService.Create(ctx, &notification_sent_history.CreateDto{
		Type:      "certificate",
		MonitorID: m.ID,
		Days:      targetDays,
	})
	if err != nil {
		return fmt.Errorf("failed to record notification sent: %w", err)
	}

	s.logger.Infof("Certificate notification recorded for monitor %s: %s", m.Name, message)
	return nil
}

func (s *ServiceImpl) isRootCertificate(certInfo *executor.CertificateInfo) bool {
	// This is a simplified check. In a real implementation, you'd maintain
	// a list of known root certificate fingerprints
	// For now, we'll consider self-signed certificates as roots
	return certInfo.Subject == certInfo.Issuer || certInfo.CertType == "ca"
}