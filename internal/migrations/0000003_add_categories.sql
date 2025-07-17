-- +goose Up
INSERT INTO categories VALUES
                      (0, 'Laptop', '0',  now(), now()),
                      (1, 'Laptop', '1',  now(), now());
-- +goose Down
-- +goose StatementBegin
DELETE FROM categories WHERE id in (0,1);
-- +goose StatementEnd