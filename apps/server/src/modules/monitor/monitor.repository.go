package monitor

import "context"

type MonitorRepository interface {
	Create(ctx context.Context, monitor *Model) (*Model, error)
	FindByID(ctx context.Context, id string) (*Model, error)
	FindAll(
		ctx context.Context,
		page int,
		limit int,
		q string,
		active *bool,
		status *int,
	) ([]*Model, error)
	FindActive(ctx context.Context) ([]*Model, error)
	UpdateFull(ctx context.Context, id string, monitor *Model) error
	UpdatePartial(ctx context.Context, id string, monitor *UpdateModel) error
	Delete(ctx context.Context, id string) error
	RemoveProxyReference(ctx context.Context, proxyId string) error
	FindByProxyId(ctx context.Context, proxyId string) ([]*Model, error)
}
