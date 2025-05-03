-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sample_table(
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sample_table;
-- +goose StatementEnd