-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN phone SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_token;
DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
