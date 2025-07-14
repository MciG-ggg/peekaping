package notification_sent_history

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, entity *CreateDto) (*Model, error)
	FindByTypeMonitorAndDays(ctx context.Context, notificationType string, monitorID string, days int) (*Model, error)
	DeleteByMonitorID(ctx context.Context, monitorID string) error
}

type ServiceImpl struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewService(
	repository Repository,
	logger *zap.SugaredLogger,
) Service {
	return &ServiceImpl{
		repository,
		logger.Named("[notification-sent-history-service]"),
	}
}

func (s *ServiceImpl) Create(ctx context.Context, entity *CreateDto) (*Model, error) {
	return s.repository.Create(ctx, entity)
}

func (s *ServiceImpl) FindByTypeMonitorAndDays(ctx context.Context, notificationType string, monitorID string, days int) (*Model, error) {
	return s.repository.FindByTypeMonitorAndDays(ctx, notificationType, monitorID, days)
}

func (s *ServiceImpl) DeleteByMonitorID(ctx context.Context, monitorID string) error {
	return s.repository.DeleteByMonitorID(ctx, monitorID)
}