-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    "id" SERIAL PRIMARY KEY,
    "phone_number" VARCHAR(191) UNIQUE NOT NULL,
    "avatar" VARCHAR(255),
    "fullname" VARCHAR(191),
    "nickname" VARCHAR(191),
    "hashed_password" VARCHAR(255) NOT NULL,
    "email" VARCHAR(191) UNIQUE,
    "province_id" BIGINT,
    "city_id" BIGINT,
    "score" INT DEFAULT 0,
    "role" user_role NOT NULL DEFAULT 'customer',
    "birth_at" TIMESTAMP,
    "is_deleted" BOOL DEFAULT false,
    "deleted_at" TIMESTAMP,
    "deleted_by" INT,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_phone_number_users ON users(phone_number);
CREATE INDEX idx_province_id_users ON users(province_id);
CREATE INDEX idx_city_id_users ON users(city_id);
CREATE INDEX idx_is_deleted_users ON users(is_deleted);

-- +migrate Down
DROP INDEX idx_phone_number_users;
DROP INDEX idx_province_id_users;
DROP INDEX idx_city_id_users;
DROP INDEX idx_is_deleted_users;
DROP TABLE IF EXISTS users;
