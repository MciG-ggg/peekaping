-- Down migration for certificate notification tracking
-- This migration removes the notification_sent_history table and indexes
-- Wrapped in a transaction for atomicity

BEGIN;

-- Drop indexes first
DROP INDEX IF EXISTS idx_notification_sent_history_sent_at;
DROP INDEX IF EXISTS idx_notification_sent_history_type_days;
DROP INDEX IF EXISTS idx_notification_sent_history_monitor_type;

-- Drop notification sent history table
DROP TABLE IF EXISTS notification_sent_history;

COMMIT;