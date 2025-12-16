-- +migrate Up
CREATE TYPE files_driver AS ENUM ('s3', 'local');

-- +migrate Down
DROP TYPE files_driver;
