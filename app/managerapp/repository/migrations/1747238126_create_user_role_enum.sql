-- +migrate Up
CREATE TYPE user_role AS ENUM ('superuser', 'support', 'story', 'file', 'notification', 'client', 'guest', 'bot', 'service');

-- +migrate Down
DROP TYPE user_role;
