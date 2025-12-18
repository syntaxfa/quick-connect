-- +migrate Up
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J0000000000000000MANAGER', 'manager-service', '$2a$14$3oa3jo6ZMDKZJ8Q45sCsM.ZARdIpgs9WwyLNPR7EHLVV8FhYtcGSC', 'manager');
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J0000000000000000000CHAT', 'chat-service', '$2a$14$jFxXyD57olvoPu9x/hJOF.htwaiKXPXiJevzvlTbr955BDs7bDkmG', 'chat');
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J000000000000000000NOTIF', 'notification-service', '$2a$14$Y5M74blYq6AVAGwJFIijFerYGqvFj1VPFJXB5e8JdkQpnZy6nzSWK', 'notification');
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J000000000000000000ADMIN', 'admin-service', '$2a$14$t3bDVpzwIDiY86C/8eIsCuL10Bub.BF0WHelgwYoWtK.Y48bxke6O', 'admin');
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J0000000000000000STORAGE', 'storage-service', '$2a$14$wPToq9WbwoV2uhm6C5A5jOVE8MJwLSr.hTCsg16/MrkqPaZxjZ6EW', 'storage');
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J000000000000000000STORY', 'story-service', '$2a$14$f7dSSVOrKFkys5dXPxfR3us4uOTvTxu799IfqGwmHN9fuV3zUbt7e', 'story');
INSERT INTO users (id, username, hashed_password, fullname) VALUES ('01J00000000000000000000RAG', 'rag-service', '$2a$14$ZDlXL9XYk5kskluWwENPpekszUqvHCugqdcE48sh8Cr3jD9/jvzmq', 'rag');

-- +migrate Down
DELETE FROM users WHERE id = '01J0000000000000000MANAGER';
DELETE FROM users WHERE id = '01J0000000000000000000CHAT';
DELETE FROM users WHERE id = '01J000000000000000000NOTIF';
DELETE FROM users WHERE id = '01J000000000000000000ADMIN';
DELETE FROM users WHERE id = '01J0000000000000000STORAGE';
DELETE FROM users WHERE id = '01J000000000000000000STORY';
DELETE FROM users WHERE id = '01J00000000000000000000RAG';
