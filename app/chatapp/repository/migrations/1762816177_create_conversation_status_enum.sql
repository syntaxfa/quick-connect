-- +migrate Up
CREATE TYPE conversation_status AS EnUM ('new', 'open', 'closed', 'bot_handling');

-- +migrate Down
DROP TYPE conversation_status;