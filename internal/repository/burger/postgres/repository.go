package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
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

func (r *Repository) SaveBurger(ctx context.Context, burgerInfo model.BurgerInfo) error {
	const fn = "repository.postgres.SaveBurger"

	stmt, err := r.db.Prepare(
		"INSERT INTO burgers (category_id, title, handle, instructions, video, data_modified) VALUES ($1, $2, $3, $4, $5, $6);",
	)
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.ExecContext(
		ctx,
		burgerInfo.CategoryId,
		burgerInfo.Title,
		burgerInfo.Handle,
		burgerInfo.Instructions,
		burgerInfo.Video,
		burgerInfo.DataModified,
	)
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

	return nil
}
