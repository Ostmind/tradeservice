-- +goose Up
INSERT INTO categories VALUES
                      (1, 'Laptop', '0',  now(), now()),
                      (2, 'Laptop', '1',  now(), now());
-- +goose Down
-- +goose StatementBegin
DELETE FROM categories WHERE id in (1,2);
-- +goose StatementEnd