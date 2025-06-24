package stats

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sqlModel struct {
	bun.BaseModel `bun:"table:stats,alias:s"`

	ID          string    `bun:"id,pk"`
	MonitorID   string    `bun:"monitor_id,notnull"`
	Timestamp   time.Time `bun:"timestamp,notnull"`
	Ping        float64   `bun:"ping,notnull,default:0"`
	PingMin     float64   `bun:"ping_min,notnull,default:0"`
	PingMax     float64   `bun:"ping_max,notnull,default:0"`
	Up          int       `bun:"up,notnull,default:0"`
	Down        int       `bun:"down,notnull,default:0"`
	Maintenance int       `bun:"maintenance,notnull,default:0"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func toDomainModelFromSQL(sm *sqlModel) *Stat {
	// Convert string IDs to ObjectIDs for compatibility with existing domain model
	objID, _ := primitive.ObjectIDFromHex(sm.ID)
	if objID.IsZero() {
		objID = primitive.NewObjectID()
	}

	monitorObjID, _ := primitive.ObjectIDFromHex(sm.MonitorID)
	if monitorObjID.IsZero() {
		monitorObjID = primitive.NewObjectID()
	}

	return &Stat{
		ID:          objID,
		MonitorID:   monitorObjID,
		Timestamp:   sm.Timestamp,
		Ping:        sm.Ping,
		PingMin:     sm.PingMin,
		PingMax:     sm.PingMax,
		Up:          sm.Up,
		Down:        sm.Down,
		Maintenance: sm.Maintenance,
	}
}

func toSQLModel(s *Stat) *sqlModel {
	return &sqlModel{
		ID:          s.ID.Hex(),
		MonitorID:   s.MonitorID.Hex(),
		Timestamp:   s.Timestamp,
		Ping:        s.Ping,
		PingMin:     s.PingMin,
		PingMax:     s.PingMax,
		Up:          s.Up,
		Down:        s.Down,
		Maintenance: s.Maintenance,
	}
}

type SQLRepositoryImpl struct {
	db *bun.DB
}

func NewSQLRepository(db *bun.DB) Repository {
	return &SQLRepositoryImpl{db: db}
}

func (r *SQLRepositoryImpl) GetOrCreateStat(ctx context.Context, monitorID primitive.ObjectID, timestamp time.Time, period StatPeriod) (*Stat, error) {
	sm := new(sqlModel)

	// Try to find existing stat
	err := r.db.NewSelect().
		Model(sm).
		Where("monitor_id = ? AND timestamp = ?", monitorID.Hex(), timestamp).
		Scan(ctx)

	if err != nil && err.Error() == "sql: no rows in result set" {
		// Create new stat if not found
		sm = &sqlModel{
			MonitorID:   monitorID.Hex(),
			Timestamp:   timestamp,
			Ping:        0,
			PingMin:     0,
			PingMax:     0,
			Up:          0,
			Down:        0,
			Maintenance: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		_, err = r.db.NewInsert().Model(sm).Returning("*").Exec(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return toDomainModelFromSQL(sm), nil
}

func (r *SQLRepositoryImpl) UpsertStat(ctx context.Context, stat *Stat, period StatPeriod) error {
	sm := toSQLModel(stat)
	sm.UpdatedAt = time.Now()

	// Use database-agnostic upsert with Bun
	_, err := r.db.NewInsert().
		Model(sm).
		On("CONFLICT (monitor_id, timestamp) DO UPDATE").
		Set("ping = EXCLUDED.ping").
		Set("ping_min = EXCLUDED.ping_min").
		Set("ping_max = EXCLUDED.ping_max").
		Set("up = EXCLUDED.up").
		Set("down = EXCLUDED.down").
		Set("maintenance = EXCLUDED.maintenance").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

func (r *SQLRepositoryImpl) FindStatsByMonitorIDAndTimeRange(ctx context.Context, monitorID primitive.ObjectID, since, until time.Time, period StatPeriod) ([]*Stat, error) {
	var sms []*sqlModel
	err := r.db.NewSelect().
		Model(&sms).
		Where("monitor_id = ? AND timestamp BETWEEN ? AND ?", monitorID.Hex(), since, until).
		Order("timestamp ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var stats []*Stat
	for _, sm := range sms {
		stats = append(stats, toDomainModelFromSQL(sm))
	}
	return stats, nil
}
