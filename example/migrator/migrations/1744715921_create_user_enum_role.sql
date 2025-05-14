-- +migrate Up
CREATE TYPE user_role AS ENUM ('admin', 'owner', 'customer');

-- +migrate Down
DROP TYPE user_role;