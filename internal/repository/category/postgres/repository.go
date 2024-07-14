package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCategory(ctx context.Context, id int) (model.Category, error) {
	const fn = "repository.postgres.GetCategory"

	var category model.Category

	stmt, err := r.db.Prepare("SELECT id, handler, title FROM categories WHERE id = $1")
	if err != nil {
		return category, fmt.Errorf("%s, %w", fn, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	err = stmt.QueryRowContext(ctx, id).Scan(&category.ID, &category.Handler, &category.Title)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return category, iErr.ErrCategoryNotFound
		}

		return category, fmt.Errorf("%s, %w", fn, err)
	}

	return category, nil
}
