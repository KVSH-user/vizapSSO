-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    phone VARCHAR UNIQUE NOT NULL,
    password_hashed VARCHAR NOT NULL,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS user_data(
    id SERIAL PRIMARY KEY,
    full_name VARCHAR,
    address VARCHAR,
    email VARCHAR,
    uid INT NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(id)
);

CREATE INDEX users_phone_idx ON users(phone);

CREATE TABLE IF NOT EXISTS refresh_token(
    token VARCHAR NOT NULL,
    uid INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS apps(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    secret VARCHAR NOT NULL UNIQUE
);

INSERT INTO apps (name, secret) VALUES ('VIZAP', 'secret123');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_token;
DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
