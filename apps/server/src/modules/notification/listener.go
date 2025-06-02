package notification

import (
	"context"
	"peekaping/src/config"
	"peekaping/src/modules/events"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/monitor_notification"
	"peekaping/src/modules/notification/integrations"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

// NotificationEventListener handles notification events
type NotificationEventListener struct {
	service                    Service
	monitorSvc                 monitor.Service
	heartbeatService           heartbeat.Service
	monitorNotificationService monitor_notification.Service
	logger                     *zap.SugaredLogger
}

type NotificationEventListenerParams struct {
	dig.In
	Service                    Service
	MonitorSvc                 monitor.Service
	HeartbeatService           heartbeat.Service
	MonitorNotificationService monitor_notification.Service
	Logger                     *zap.SugaredLogger
	Config                     *config.Config
}

func NewNotificationEventListener(p NotificationEventListenerParams) *NotificationEventListener {
	RegisterNotifier("smtp", integrations.NewEmailSender(p.Logger))
	RegisterNotifier("telegram", integrations.NewTelegramSender(p.Logger))
	RegisterNotifier("webhook", integrations.NewWebhookSender(p.Logger))
	RegisterNotifier("slack", integrations.NewSlackSender(p.Logger, p.Config))

	return &NotificationEventListener{
		service:                    p.Service,
		monitorSvc:                 p.MonitorSvc,
		heartbeatService:           p.HeartbeatService,
		monitorNotificationService: p.MonitorNotificationService,
		logger:                     p.Logger,
	}
}

// Subscribe subscribes to NotifyEvent and sends notifications
func (l *NotificationEventListener) Subscribe(eventBus *events.EventBus) {
	eventBus.Subscribe(events.MonitorStatusChanged, l.handleNotifyEvent)
}

func (l *NotificationEventListener) handleNotifyEvent(event events.Event) {
	ctx := context.Background()

	hb, ok := event.Payload.(*heartbeat.Model)
	if !ok {
		l.logger.Errorf("Invalid handleNotifyEvent event payload type: %v", event.Payload)
		return
	}

	monitorID := hb.MonitorID

	l.logger.Infof("Notification event received for monitor: %s", monitorID)

	// Get monitor-notification records
	monitorNotifications, err := l.monitorNotificationService.FindByMonitorID(ctx, monitorID)
	if err != nil {
		l.logger.Errorf("Failed to get monitor-notification records: %v", err)
		return
	}

	var notifications []*Model
	for _, mn := range monitorNotifications {
		l.logger.Infof("Monitor notification: %s", mn.NotificationID)
		notification, err := l.service.FindByID(ctx, mn.NotificationID)
		if err != nil {
			l.logger.Errorf("Failed to get notification by ID: %s, error: %v", mn.NotificationID, err)
			continue
		}
		if notification != nil {
			notifications = append(notifications, notification)
		} else {
			l.logger.Warnf("Notification not found for monitor-notification: %s", mn.NotificationID)
		}
	}

	// Fetch monitor details for context
	monitorModel, err := l.monitorSvc.FindByID(ctx, monitorID)
	if err != nil || monitorModel == nil {
		l.logger.Warn("Monitor not found for notification context")
		return
	}

	for _, notification := range notifications {
		integration, ok := GetNotifier(notification.Type)
		if !ok {
			l.logger.Warnf("No integration registered for notification type: %s", notification.Type)
			continue
		}
		if notification.Config == nil {
			l.logger.Warnf("No config for notification: %s", notification.Name)
			continue
		}

		// validate config
		if err := integration.Validate(*notification.Config); err != nil {
			l.logger.Errorf("Failed to validate notification config: %s, error: %v", notification.Name, err)
			continue
		}

		err := integration.Send(ctx, *notification.Config, hb.Msg, monitorModel, hb)
		if err != nil {
			l.logger.Errorf("Failed to send notification: %s, error: %v", notification.Name, err)
		} else {
			l.logger.Infof("Notification sent to: %s for monitor: %s", notification.Name, monitorID)
		}
	}
}
