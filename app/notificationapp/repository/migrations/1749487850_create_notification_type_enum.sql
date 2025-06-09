-- +migrate UP
CREATE TYPE notification_type AS ENUM ('optional', 'info', 'promotion', 'critical');

-- +migrate Down
DROP TYPE notification_type;