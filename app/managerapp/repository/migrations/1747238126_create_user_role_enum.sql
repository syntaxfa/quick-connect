-- +migrate Up
CREATE TYPE user_role AS ENUM ('superuser', 'admin');

-- +migrate Down
DROP TYPE user_role;