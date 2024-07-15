-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories
(
    id SERIAL PRIMARY KEY,
    handle varchar(255) NOT NULL UNIQUE,
    title varchar(255) NOT NULL UNIQUE
);
INSERT INTO categories (handle, title) VALUES ('other-unknown', 'Other/Unknown')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories
-- +goose StatementEnd
