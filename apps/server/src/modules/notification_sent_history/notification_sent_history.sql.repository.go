package notification_sent_history

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type sqlModel struct {
	bun.BaseModel `bun:"table:notification_sent_history,alias:nsh"`

	ID        string    `bun:"id,pk"`
	Type      string    `bun:"type,notnull"`
	MonitorID string    `bun:"monitor_id,notnull"`
	Days      int       `bun:"days,notnull"`
	SentAt    time.Time `bun:"sent_at,nullzero,notnull,default:current_timestamp"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func toDomainModelFromSQL(sm *sqlModel) *Model {
	return &Model{
		ID:        sm.ID,
		Type:      sm.Type,
		MonitorID: sm.MonitorID,
		Days:      sm.Days,
		SentAt:    sm.SentAt,
		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}
}

func toSQLModel(dto *CreateDto) *sqlModel {
	now := time.Now()
	return &sqlModel{
		ID:        uuid.New().String(),
		Type:      dto.Type,
		MonitorID: dto.MonitorID,
		Days:      dto.Days,
		SentAt:    now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type SQLRepositoryImpl struct {
	db *bun.DB
}

func NewSQLRepository(db *bun.DB) Repository {
	return &SQLRepositoryImpl{db: db}
}

func (r *SQLRepositoryImpl) Create(ctx context.Context, entity *CreateDto) (*Model, error) {
	sm := toSQLModel(entity)

	_, err := r.db.NewInsert().Model(sm).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainModelFromSQL(sm), nil
}

func (r *SQLRepositoryImpl) FindByTypeMonitorAndDays(ctx context.Context, notificationType string, monitorID string, days int) (*Model, error) {
	sm := new(sqlModel)
	err := r.db.NewSelect().
		Model(sm).
		Where("type = ? AND monitor_id = ? AND days <= ?", notificationType, monitorID, days).
		Order("days DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return toDomainModelFromSQL(sm), nil
}

func (r *SQLRepositoryImpl) DeleteByMonitorID(ctx context.Context, monitorID string) error {
	_, err := r.db.NewDelete().Model((*sqlModel)(nil)).Where("monitor_id = ?", monitorID).Exec(ctx)
	return err
}