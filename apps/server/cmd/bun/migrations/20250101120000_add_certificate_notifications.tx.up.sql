-- Add notification history tracking for certificate expiry notifications
-- This migration adds support for tracking sent notifications to prevent duplicates
-- Wrapped in a transaction for atomicity

-- Notification sent history table
CREATE TABLE IF NOT EXISTS notification_sent_history (
    id UUID PRIMARY KEY,
    type VARCHAR(50) NOT NULL, -- 'certificate', 'monitor_status', etc.
    monitor_id UUID NOT NULL,
    days INTEGER NOT NULL, -- notification threshold days (e.g., 7, 14, 21)
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_notification_sent_history_monitor_type ON notification_sent_history(monitor_id, type);
CREATE INDEX IF NOT EXISTS idx_notification_sent_history_type_days ON notification_sent_history(type, monitor_id, days);
CREATE INDEX IF NOT EXISTS idx_notification_sent_history_sent_at ON notification_sent_history(sent_at);