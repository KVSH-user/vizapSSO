-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     phone VARCHAR(15) NOT NULL UNIQUE,
                                     password_hashed TEXT NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users_data (
                                          id SERIAL PRIMARY KEY,
                                          full_name VARCHAR(255) NOT NULL,
                                          email VARCHAR(255) NOT NULL,
                                          user_id INT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS addresses (
                                         id SERIAL PRIMARY KEY,
                                         full_address TEXT NOT NULL,
                                         is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
                                         user_id INT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              id SERIAL PRIMARY KEY,
                                              token TEXT NOT NULL,
                                              created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                              is_valid BOOLEAN NOT NULL DEFAULT TRUE,
                                              user_id INT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS apps(
                                   id SERIAL PRIMARY KEY,
                                   name VARCHAR NOT NULL UNIQUE,
                                   secret VARCHAR NOT NULL UNIQUE
);

INSERT INTO apps (name, secret) VALUES ('VIZAP', 'secret');

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_data_user_id ON users_data(user_id);
CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS users_data;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_data_user_id;
DROP INDEX IF EXISTS idx_addresses_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
-- +goose StatementEnd
