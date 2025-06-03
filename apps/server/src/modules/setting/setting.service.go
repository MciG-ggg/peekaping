package setting

import "context"

type Service interface {
	GetByKey(ctx context.Context, key string) (*Model, error)
	SetByKey(ctx context.Context, key string, entity *CreateUpdateDto) (*Model, error)
	DeleteByKey(ctx context.Context, key string) error
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

func (mr *ServiceImpl) GetByKey(ctx context.Context, key string) (*Model, error) {
	return mr.repository.GetByKey(ctx, key)
}

func (mr *ServiceImpl) SetByKey(ctx context.Context, key string, entity *CreateUpdateDto) (*Model, error) {
	return mr.repository.SetByKey(ctx, key, entity)
}

func (mr *ServiceImpl) DeleteByKey(ctx context.Context, key string) error {
	return mr.repository.DeleteByKey(ctx, key)
}
