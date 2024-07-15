package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/ingredient/model"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllIngredients(ctx context.Context, params serviceModel.FetchParam) ([]model.Ingredient, error) {
	const fn = "repository.postgres.GetAllIngredient"

	var ingredients = make([]model.Ingredient, 0)

	queryBuilder := sq.Select("id", "handle", "title", "description").
		From("ingredients").PlaceholderFormat(sq.Dollar)

	if params.Title != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{
			"title": params.Title,
		})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return ingredients, fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return ingredients, fmt.Errorf("%s, %w", fn, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return ingredients, fmt.Errorf("%s, %w", fn, err)
	}
	for rows.Next() {
		var ingredient model.Ingredient
		if err := rows.Scan(
			&ingredient.ID,
			&ingredient.Title,
			&ingredient.Handle,
			&ingredient.Description,
		); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (r *Repository) GetIngredient(ctx context.Context, id int) (model.Ingredient, error) {
	const fn = "repository.postgres.GetIngredient"

	var ingredient model.Ingredient

	stmt, err := r.db.Prepare("SELECT id, handle, title, description FROM ingredients WHERE id = $1")
	if err != nil {
		return ingredient, fmt.Errorf("%s, %w", fn, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	err = stmt.QueryRowContext(ctx, id).Scan(
		&ingredient.ID,
		&ingredient.Handle,
		&ingredient.Title,
		&ingredient.Description,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ingredient, iErr.ErrIngredientNotFound
		}

		return ingredient, fmt.Errorf("%s, %w", fn, err)
	}

	return ingredient, nil
}
