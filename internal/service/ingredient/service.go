package ingredient

import (
	"context"
	"errors"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/ingredient/model"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

type Repository interface {
	GetIngredient(ctx context.Context, id int) (model.Ingredient, error)
	GetAllIngredients(ctx context.Context, params serviceModel.FetchParam) ([]model.Ingredient, error)
}

type Service struct {
	logger               *slog.Logger
	ingredientRepository Repository
}

func New(repo Repository, logger *slog.Logger) *Service {
	return &Service{ingredientRepository: repo, logger: logger}
}

func (s *Service) GetAllIngredients(ctx context.Context, params serviceModel.FetchParam) ([]model.Ingredient, uint64, error) {
	const fn = "service.ingredient.GetAllIngredients"

	var nextCursor uint64

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
		"params", params,
	)

	ingredients, err := s.ingredientRepository.GetAllIngredients(ctx, params)
	if err != nil {
		log.Info("fail to get ingredient")

		return nil, nextCursor, err
	}

	if len(ingredients) > 0 {
		nextCursor = uint64(ingredients[len(ingredients)-1].ID)
	}

	return ingredients, nextCursor, nil
}

func (s *Service) GetIngredient(ctx context.Context, id int) (model.Ingredient, error) {
	const fn = "service.ingredient.Get"

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	ingredient, err := s.ingredientRepository.GetIngredient(ctx, id)
	if errors.Is(err, iErr.ErrIngredientNotFound) {
		log.Info("ingredient not found", "id", id)

		return model.Ingredient{}, err
	}

	if err != nil {
		log.Info("fail to get ingredient", "id", id)

		return model.Ingredient{}, err
	}

	return ingredient, nil
}
