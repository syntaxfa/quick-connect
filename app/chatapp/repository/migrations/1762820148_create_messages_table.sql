-- +migrate Up
CREATE TABLE IF NOT EXISTS messages (
    "id" VARCHAR(26) PRIMARY KEY,
    "conversation_id" VARCHAR(26) NOT NULL,
    "sender_id" VARCHAR(26) NOT NULL,
    "message_type" message_type NOT NULL DEFAULT 'text',
    "content" TEXT NULL,
    "metadata" JSONB NULL,
    "replied_to_message_id" VARCHAR(26) NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "read_at" TIMESTAMPTZ NULL,
    FOREIGN KEY ("conversation_id") REFERENCES conversations("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages("conversation_id");
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages("sender_id");
CREATE INDEX IF NOT EXISTS idx_messages_replied_to_message_id ON messages("replied_to_message_id");

-- +migrate Down
DROP INDEX IF EXISTS idx_messages_replied_to_message_id;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_conversation_id;
DROP TABLE IF EXISTS messages;