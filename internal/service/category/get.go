package category

import (
	"context"
	"errors"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	"github.com/go-chi/chi/v5/middleware"
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

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	category, err := s.categoryRepository.GetCategory(ctx, id)
	if errors.Is(err, iErr.ErrCategoryNotFound) {
		log.Info("category not found", "id", id)

		return model.Category{}, err
	}

	if err != nil {
		log.Info("fail to get category", "id", id)

		return model.Category{}, err
	}

	return category, nil
}
