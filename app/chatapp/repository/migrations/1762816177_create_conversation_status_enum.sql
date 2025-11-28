-- +migrate Up
CREATE TYPE conversation_status AS EKHAM ("new", 'open', 'closed', 'bot_handling');

-- +migrate Down
DROP TYPE conversation_status;