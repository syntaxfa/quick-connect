-- +migrate Up
CREATE TABLE IF NOT EXISTS files (
    "id" VARCHAR(26) PRIMARY KEY,
    "uploader_id" VARCHAR(26) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "key" VARCHAR(512) NOT NULL,
    "mime_type" VARCHAR(100) NOT NULL,
    "size" BIGINT NOT NULL,
    "driver" files_driver NOT NULL,
    "bucket" VARCHAR(100) NULL,
    "is_public" BOOLEAN DEFAULT false,
    "is_confirmed" BOOLEAN DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ NULL
);
CREATE INDEX idx_files_id ON files(id);
CREATE INDEX idx_files_uploader_id ON files(uploader_id);
CREATE INDEX idx_files_deleted_at ON files(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_files_is_confirmed ON files(is_confirmed) WHERE is_confirmed IS false;

-- +migrate Down
DROP INDEX IF EXISTS idx_files_is_confirmed;
DROP INDEX IF EXISTS idx_files_deleted_at;
DROP INDEX IF EXISTS idx_files_uploader_id;
DROP INDEX IF EXISTS idx_files_id;
DROP TABLE IF EXISTS files;
