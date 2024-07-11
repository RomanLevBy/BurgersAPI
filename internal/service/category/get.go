package category

import (
	"context"
	"errors"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	"log/slog"
)

type Repository interface {
	GetCategory(ctx context.Context, id int) (model.Category, error)
}

type Service struct {
	logger             *slog.Logger
	categoryRepository Repository
}

func New(repo Repository, logger *slog.Logger) *Service {
	return &Service{categoryRepository: repo, logger: logger}
}

func (s *Service) GetCategory(ctx context.Context, id int) (model.Category, error) {
	const fn = "service.category.Get"

	s.logger.With(slog.String("fn", fn))

	category, err := s.categoryRepository.GetCategory(ctx, id)
	if errors.Is(err, iErr.ErrCategoryNotFound) {
		s.logger.Info("category not found", "id", id)

		return model.Category{}, err
	}

	if err != nil {
		s.logger.Info("Fail to get category", "id", id)

		return model.Category{}, err
	}

	return category, nil
}
