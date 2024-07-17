package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/burger/model"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllBurgers(ctx context.Context, params serviceModel.FetchParam) ([]model.Burger, error) {
	const fn = "repository.postgres.GetAllBurger"

	var burgers = make([]model.Burger, 0)

	queryBuilder := sq.
		Select(
			"burgers.id",
			"burgers.handle",
			"burgers.title",
			"instructions",
			"video",
			"data_modified",
		).
		From("burgers").
		PlaceholderFormat(sq.Dollar)

	if params.Limit > 0 {
		queryBuilder = queryBuilder.Limit(params.Limit)
	}

	if params.CursorID > 0 {
		queryBuilder = queryBuilder.Where(sq.Gt{
			"id": params.CursorID,
		})
	}

	if params.Title != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{
			"title": params.Title,
		})
	}

	if params.TitlePath != "" {
		queryBuilder = queryBuilder.Where("burgers.title LIKE ?", fmt.Sprintf("%s%%", params.TitlePath))
	}

	if params.TitlePath == "" {
		queryBuilder = queryBuilder.Join("burgers_ingredients as bi ON burgers.id = bi.burger_id")
		queryBuilder = queryBuilder.Join("ingredients ON bi.ingredient_id = ingredients.id")
		queryBuilder = queryBuilder.Where("ingredients.title LIKE ?", fmt.Sprintf("%s%%", "test"))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return burgers, fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return burgers, fmt.Errorf("%s, %w", fn, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return burgers, fmt.Errorf("%s, %w", fn, err)
	}
	for rows.Next() {
		var burger model.Burger
		if err := rows.Scan(
			&burger.ID,
			&burger.Title,
			&burger.Handle,
			&burger.Instructions,
			&burger.Video,
			&burger.DataModified,
		); err != nil {
			return nil, err
		}
		burgers = append(burgers, burger)
	}

	return burgers, nil
}

func (r *Repository) GetBurger(ctx context.Context, id int) (model.Burger, error) {
	const fn = "repository.postgres.GetBurger"

	var burger model.Burger

	stmt, err := r.db.Prepare("SELECT id, handle, title, instructions, video, data_modified FROM burgers WHERE id = $1")
	if err != nil {
		return burger, fmt.Errorf("%s, %w", fn, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	err = stmt.QueryRowContext(ctx, id).Scan(
		&burger.ID,
		&burger.Handle,
		&burger.Title,
		&burger.Instructions,
		&burger.Video,
		&burger.DataModified,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return burger, iErr.ErrBurgerNotFound
		}

		return burger, fmt.Errorf("%s, %w", fn, err)
	}

	return burger, nil
}

func (r *Repository) SaveBurger(ctx context.Context, burgerInfo serviceModel.BurgerInfo) error {
	const fn = "repository.postgres.SaveBurger"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}
	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
		}
	}(tx)

	stmt, err := tx.Prepare(
		"INSERT INTO burgers (category_id, title, handle, instructions, video, data_modified) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
	)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	var lastInsertId int64
	err = stmt.QueryRowContext(
		ctx,
		burgerInfo.CategoryId,
		burgerInfo.Title,
		burgerInfo.Handle,
		burgerInfo.Instructions,
		burgerInfo.Video,
		burgerInfo.DataModified,
	).Scan(&lastInsertId)
	if err != nil {
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			if pgErr.Code == "23505" {
				return iErr.ErrBurgerExists
			}
		}

		return fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err = tx.Prepare(
		"INSERT INTO burgers_ingredients (burger_id, ingredient_id, instruction) VALUES ($1, $2, $3);",
	)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	for _, ingredient := range burgerInfo.Ingredients {
		_, err = stmt.ExecContext(
			ctx,
			lastInsertId,
			ingredient.IngredientId,
			ingredient.Instruction,
		)
		if err != nil {
			return fmt.Errorf("%s, %w", fn, err)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	return nil
}
