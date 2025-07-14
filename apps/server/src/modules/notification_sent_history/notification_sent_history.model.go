package notification_sent_history

import "time"

type Model struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`        // 'certificate', 'monitor_status', etc.
	MonitorID string    `json:"monitor_id"`
	Days      int       `json:"days"`        // notification threshold days (e.g., 7, 14, 21)
	SentAt    time.Time `json:"sent_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDto struct {
	Type      string `json:"type" validate:"required"`
	MonitorID string `json:"monitor_id" validate:"required"`
	Days      int    `json:"days" validate:"required,min=1"`
}