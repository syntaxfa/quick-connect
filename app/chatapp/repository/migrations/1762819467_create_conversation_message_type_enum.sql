-- +migrate Up
CREATE TYPE message_type AS ENUM ('text', 'media', 'system');

-- +migrate Down
DROP TYPE message_type;
