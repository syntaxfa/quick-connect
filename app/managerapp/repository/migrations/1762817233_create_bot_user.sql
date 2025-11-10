-- +migrate Up
INSERT INTO users (id, username, fullname)
VALUES ('01J00000000000000000000BOT', 'bot', 'Quick Connect AI')
ON CONFLICT (id) DO NOTHING;

-- +migrate Down
DELETE FROM users WHERE id = '01J00000000000000000000BOT';