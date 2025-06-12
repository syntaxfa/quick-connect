-- +migrate Up
CREATE TYPE notification_type AS ENUM ('optional', 'info', 'promotion', 'critical', 'direct');

-- +migrate Down
DROP TYPE notification_type;