-- +goose Up
INSERT INTO products VALUES
                      (1, 'Lenovo',now(), now()),
                      (2, 'Macbook',now(), now());
-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE id in (1,2);
-- +goose StatementEnd