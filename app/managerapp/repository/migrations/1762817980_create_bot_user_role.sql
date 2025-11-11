-- +migrate Up
INSERT INTO user_roles (user_id, role)
VALUES ('01J00000000000000000000BOT', 'bot')
ON CONFLICT (user_id, role) DO NOTHING;

-- +migrate Down
DELETE FROM user_roles WHERE user_id = '01J00000000000000000000BOT';