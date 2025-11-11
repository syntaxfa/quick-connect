-- +migrate Up
CREATE TABLE IF NOT EXISTS conversations (
    "id" VARCHAR(26) PRIMARY KEY,
    "client_user_id" VARCHAR(26) NOT NULL,
    "assigned_support_id" VARCHAR(26) NULL,
    "status" conversation_status NOT NULL,
    "last_message_snippet" TEXT NULL,
    "last_message_sender_id" VARCHAR(26) NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "closed_at" TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_conversation_client_user_id ON conversations("client_user_id");
CREATE INDEX IF NOT EXISTS idx_conversations_assigned_support_id ON conversations("assigned_support_id");
CREATE INDEX IF NOT EXISTS idx_conversations_status ON conversations("status");

-- +migrate Down
DROP INDEX IF EXISTS idx_conversations_status;
DROP INDEX IF EXISTS idx_conversations_assigned_support_id;
DROP INDEX IF EXISTS idx_conversation_client_user_id;
DROP TABLE IF EXISTS conversations;