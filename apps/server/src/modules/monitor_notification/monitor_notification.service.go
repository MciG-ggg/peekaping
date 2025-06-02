package monitor_notification

import "context"

type Service interface {
	Create(ctx context.Context, monitorID string, notificationID string) (*Model, error)
	FindByID(ctx context.Context, id string) (*Model, error)
	Delete(ctx context.Context, id string) error
	FindByMonitorID(ctx context.Context, monitorID string) ([]*Model, error)
	DeleteByMonitorID(ctx context.Context, monitorID string) error
	DeleteByNotificationID(ctx context.Context, notificationID string) error
}

type ServiceImpl struct {
	repository Repository
}

func NewService(
	repository Repository,
) Service {
	return &ServiceImpl{
		repository,
	}
}

func (mr *ServiceImpl) Create(ctx context.Context, monitorID string, notificationID string) (*Model, error) {
	createModel := &Model{
		MonitorID:      monitorID,
		NotificationID: notificationID,
	}

	return mr.repository.Create(ctx, createModel)
}

func (mr *ServiceImpl) FindByID(ctx context.Context, id string) (*Model, error) {
	return mr.repository.FindByID(ctx, id)
}

func (mr *ServiceImpl) Delete(ctx context.Context, id string) error {
	return mr.repository.Delete(ctx, id)
}

func (mr *ServiceImpl) FindByMonitorID(ctx context.Context, monitorID string) ([]*Model, error) {
	return mr.repository.FindByMonitorID(ctx, monitorID)
}

func (mr *ServiceImpl) DeleteByMonitorID(ctx context.Context, monitorID string) error {
	return mr.repository.DeleteByMonitorID(ctx, monitorID)
}

func (mr *ServiceImpl) DeleteByNotificationID(ctx context.Context, notificationID string) error {
	return mr.repository.DeleteByNotificationID(ctx, notificationID)
}
