-- +goose Up
-- +goose StatementBegin
ALTER TABLE refresh_token
ADD COLUMN is_active BOOL NOT NULL DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_token;
DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
