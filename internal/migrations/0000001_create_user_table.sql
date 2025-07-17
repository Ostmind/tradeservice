-- +goose Up
CREATE TABLE users (
                       id INTEGER PRIMARY KEY,
                       username TEXT,
                       name TEXT,
                       surname TEXT
);

-- +goose Down
DROP TABLE users;