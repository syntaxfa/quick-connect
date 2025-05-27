-- +migrate Up
CREATE TABLE IF NOT EXISTS files_access (
    "id" SERIAL PRIMARY KEY,
    "permission" file_permissions[] NOT NULL,
    "file_id" INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    "client_id" VARCHAR(10) NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS files_access;