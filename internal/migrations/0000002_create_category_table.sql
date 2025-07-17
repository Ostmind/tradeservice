-- +goose Up
CREATE TABLE categories (
                       id INTEGER PRIMARY KEY,
                       name TEXT,
                       product_id TEXT,
                       created_at date,
                       updated_at date
);

-- +goose Down
DROP TABLE categories;