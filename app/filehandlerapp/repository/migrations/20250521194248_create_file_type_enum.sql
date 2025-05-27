-- +migrate Up
CREATE TYPE file_types AS ENUM ('chat');

-- +migrate Down
DROP TYPE file_types;
