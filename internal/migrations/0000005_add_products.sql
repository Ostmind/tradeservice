-- +goose Up
INSERT INTO products VALUES
                      (0, 'Lenovo',now(), now()),
                      (1, 'Macbook',now(), now());
-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE id in (0,1);
-- +goose StatementEnd