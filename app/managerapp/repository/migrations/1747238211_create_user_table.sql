-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(191) UNIQUE NOT NULL,
    "hashed_password" VARCHAR(255),
    "fullname" VARCHAR(191) NULL,
    "avatar" VARCHAR(255) NULL,
    "role" user_role NOT NULL DEFAULT 'admin',
    "last_online_at" TIMESTAMP DEFAULT NOW(),
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_username_users ON users(username);

-- +migrate Down
DROP INDEX idx_username_users;
DROP TABLE IF EXISTS users;
