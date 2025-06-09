-- +migrate Up
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'failed', 'retrying', 'ignored');

-- +migrate Down
DROP TYPE notification_status;