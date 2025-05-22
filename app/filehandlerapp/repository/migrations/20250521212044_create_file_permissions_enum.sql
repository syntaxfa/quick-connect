-- +migrate Up
CREATE TYPE file_permissions AS ENUM ('r', 'w', 'u','d');

-- +migrate Down
DROP TYPE file_permissions;
