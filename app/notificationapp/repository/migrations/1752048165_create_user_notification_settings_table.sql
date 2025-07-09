-- +migrate Up
CREATE TABLE IF NOT EXISTS user_notification_settings (
    "id" VARCHAR(26) PRIMARY KEY,
    "user_id" VARCHAR(26) NOT NULL UNIQUE,
    "lang" VARCHAR(26) NOT NULL,
    "ignore_channels" JSONB NULL
);
CREATE INDEX idx_user_id_user_notification_settings ON user_notification_settings(user_id);

-- +migrate Down
DROP INDEX idx_user_id_user_notification_settings;
DROP TABLE IF EXISTS user_notification_settings;
