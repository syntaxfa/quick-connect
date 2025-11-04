-- +migrate Up
CREATE TABLE IF NOT EXISTS external_users (
    "user_id" VARCHAR(26) NOT NULL UNIQUE,
    "external_user_id" VARCHAR(255) NOT NULL UNIQUE
);
CREATE INDEX idx_external_user_id_external_user ON external_users(external_user_id);

-- +migrate Down
DROP INDEX idx_external_user_id_external_user;
DROP TABLE IF EXISTS external_users;