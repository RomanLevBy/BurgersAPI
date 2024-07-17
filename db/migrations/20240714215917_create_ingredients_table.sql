-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ingredients
(
    id SERIAL PRIMARY KEY,
    handle varchar(255) NOT NULL UNIQUE,
    title varchar(255) NOT NULL UNIQUE,
    description TEXT NOT NULL
);
CREATE INDEX ingredients_title_index ON ingredients(title varchar_pattern_ops);
INSERT INTO ingredients (handle, title, description) VALUES ('other-unknown', 'Other/Unknown', 'Other/Unknown');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ingredients
-- +goose StatementEnd
