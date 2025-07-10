-- +migrate Up
CREATE TABLE IF NOT EXISTS notifications (
    "id" VARCHAR(26) PRIMARY KEY,
    "user_id" VARCHAR(26) NOT NULL,
    "type" notification_type NOT NULL,
    "data" JSONB NULL,
    "template_name" VARCHAR(255) NOT NULL,
    "dynamic_body_data" JSONB NULL,
    "dynamic_title_data" JSONB NULL,
    "is_read" BOOL DEFAULT false,
    "is_in_app" BOOL DEFAULT false,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "overall_status" notification_status NOT NULL DEFAULT 'pending',
    "channel_deliveries" JSONB NOT NULL DEFAULT '[]'::jsonb
);
CREATE INDEX idx_user_id_notification ON notifications(user_id);
CREATE INDEX idx_type_notification ON notifications(type);
CREATE INDEX idx_is_read_notification ON notifications(is_read);
CREATE INDEX idx_is_in_app_notification ON notifications(is_in_app);
CREATE INDEX idx_created_at_notification ON notifications(created_at);
CREATE INDEX idx_overall_status_notifications ON notifications(overall_status);
CREATE INDEX idx_template_name_notifications ON notifications(template_name);

-- +migrate Down
DROP INDEX idx_user_id_notification;
DROP INDEX idx_type_notification;
DROP INDEX idx_is_read_notification;
DROP INDEX idx_is_in_app_notification;
DROP INDEX idx_created_at_notification;
DROP INDEX idx_overall_status_notifications;
DROP INDEX idx_template_name_notifications;
DROP TABLE IF EXISTS notifications;
