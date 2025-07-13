-- +migrate Up
CREATE TABLE IF NOT EXISTS templates (
    "id" VARCHAR(26) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "contents" JSONB NOT NULL,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_name_templates ON templates(name);
CREATE INDEX idx_created_at_templates ON templates(created_at);

-- +migrate Down
DROP INDEX idx_name_templates;
DROP INDEX idx_created_at_templates;
DROP TABLE IF EXISTS templates;
