-- +migrate Up
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'failed', 'retrying', 'ignored', 'mixed');

-- +migrate Down
DROP TYPE notification_status;