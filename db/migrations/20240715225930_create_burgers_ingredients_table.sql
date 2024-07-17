-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS burgers_ingredients
(
    id SERIAL PRIMARY KEY,
    burger_id INT,
    ingredient_id INT,
    instruction TEXT,
    CONSTRAINT fk_burger
        FOREIGN KEY(burger_id)
            REFERENCES burgers(id) on DELETE CASCADE,
    CONSTRAINT fk_ingredient
        FOREIGN KEY(ingredient_id)
            REFERENCES ingredients(id) on DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE burgers_ingredients;
-- +goose StatementEnd
