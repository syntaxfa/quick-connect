-- +migrate Up
CREATE TABLE IF NOT EXISTS user_roles (
    "user_id" VARCHAR(26) NOT NULL,
    "role" user_role NOT NULL,
    PRIMARY KEY ("user_id", "role"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS user_roles;