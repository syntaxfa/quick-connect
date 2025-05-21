-- +migrate Up
CREATE TABLE IF NOT EXISTS files (
    "id" SERIAL PRIMARY KEY,
    "type" file_types NOT NULL,
    "type_id" INTEGER NOT NULL,
    "extension" VARCHAR(10) NOT NULL,
    "storage_type" storage_type NOT NULL,
    "size" INTEGER NOT NULL,
    "content_type" VARCHAR(50) NOT NULL,
    "is_public" BOOLEAN NOT NULL,
    "is_deleted" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW(),
);

-- +migrate Down
DROP TABLE IF EXISTS files;