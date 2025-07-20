-- +migrate Up
CREATE TABLE IF NOT EXISTS subscriber_profiles (
    "id" VARCHAR(26) PRIMARY KEY,
    "user_id" VARCHAR(26) NOT NULL,
    "email" VARCHAR(120) NULL,
    "phone_number" VARCHAR(20) NULL,
    "push_tokens" TEXT[],
    "created_at" TIMESTAMP default NOW(),
    "updated_at" TIMESTAMP default NOW()
);
CREATE INDEX idx_user_id_subscriber_profiles ON subscriber_profiles(user_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_user_id_subscriber_profiles;
DROP TABLE IF EXISTS subscriber_profiles;
