-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories
(
    id SERIAL PRIMARY KEY,
    handler varchar(255) NOT NULL UNIQUE,
    title varchar(255) NOT NULL UNIQUE
);
INSERT INTO categories (handler, title) VALUES ('other-unknown', 'Other/Unknown')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories
-- +goose StatementEnd
