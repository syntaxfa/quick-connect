-- +migrate Up
CREATE TABLE IF NOT EXISTS stories (
    "id" VARCHAR(26) PRIMARY KEY,
    "media_file_id" VARCHAR(26) NOT NULL,
    "title" VARCHAR(255) NULL,
    "caption" TEXT NULL,
    "link_url" VARCHAR(555) NULL,
    "link_text" VARCHAR(100) NULL,
    "duration_seconds" INT NOT NULL,
    "is_active" BOOLEAN DEFAULT true,
    "publish_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "expires_at" TIMESTAMPTZ NOT NULL,
    "view_count" BIGINT DEFAULT 0,
    "creator_id" VARCHAR(26) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_stories_feed_schedule ON stories(is_active, publish_at, expires_at);
CREATE INDEX IF NOT EXISTS idx_stories_creator ON stories(creator_id, created_at DESC);

-- +migrate Down
DROP INDEX IF EXISTS idx_stories_creator;
DROP INDEX IF EXISTS idx_stories_feed_schedule;
DROP TABLE IF EXISTS stories;