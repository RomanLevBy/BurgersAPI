-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS burgers
(
    id SERIAL PRIMARY KEY,
    category_id INT,
    handle varchar(255) NOT NULL UNIQUE,
    title varchar(255) NOT NULL UNIQUE,
    instructions TEXT NOT NULL,
    video varchar(1000),
    data_modified TIMESTAMP,
    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
            REFERENCES categories(id)
);
CREATE INDEX burgers_title_index ON burgers(title varchar_pattern_ops);
-- INSERT INTO burgers
-- SELECT i, 1, md5(random()::text), md5(random()::text), md5(random()::text), md5(random()::text), now()
-- FROM generate_series(1, 700000) AS i;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE burgers
-- +goose StatementEnd
