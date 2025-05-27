-- +migrate Up
CREATE TYPE storage_types AS ENUM ('local', 's3');

-- +migrate Down
DROP TYPE storage_types;
