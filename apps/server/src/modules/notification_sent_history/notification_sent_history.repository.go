package notification_sent_history

import "context"

type Repository interface {
	Create(ctx context.Context, entity *CreateDto) (*Model, error)
	FindByTypeMonitorAndDays(ctx context.Context, notificationType string, monitorID string, days int) (*Model, error)
	DeleteByMonitorID(ctx context.Context, monitorID string) error
}