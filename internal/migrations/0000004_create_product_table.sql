-- +goose Up
CREATE TABLE products (
                       id INTEGER PRIMARY KEY,
                       name TEXT,
                       created_at date,
                       updated_at date
);

-- +goose Down
DROP TABLE products;