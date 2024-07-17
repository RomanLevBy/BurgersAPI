package burger

import (
	"context"
	"errors"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/lib/logger/sl"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/burger/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iancoleman/strcase"
	"log/slog"
	"time"
)

type Repository interface {
	GetBurger(ctx context.Context, id int) (model.Burger, error)
	SaveBurger(ctx context.Context, burgerInfo serviceModel.BurgerInfo) error
	GetAllBurgers(ctx context.Context, params serviceModel.FetchParam) ([]model.Burger, error)
}

type IngredientServer interface {
	GetIngredient(ctx context.Context, id int) (model.Ingredient, error)
}

type Service struct {
	logger           *slog.Logger
	ingredientServer IngredientServer
	burgerRepository Repository
}

func New(repo Repository, ingredientServer IngredientServer, logger *slog.Logger) *Service {
	return &Service{burgerRepository: repo, ingredientServer: ingredientServer, logger: logger}
}

func (s *Service) GetAllBurgers(ctx context.Context, params serviceModel.FetchParam) ([]model.Burger, uint64, error) {
	const fn = "service.burger.GetAllBurgers"

	var nextCursor uint64

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
		"params", params,
	)

	burgers, err := s.burgerRepository.GetAllBurgers(ctx, params)
	if err != nil {
		log.Info("fail to get burger")

		return nil, nextCursor, err
	}

	if len(burgers) > 0 {
		nextCursor = uint64(burgers[len(burgers)-1].ID)
	}

	return burgers, nextCursor, nil
}

func (s *Service) GetBurger(ctx context.Context, id int) (model.Burger, error) {
	const fn = "service.burger.Get"

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	burger, err := s.burgerRepository.GetBurger(ctx, id)
	if errors.Is(err, iErr.ErrBurgerNotFound) {
		log.Info("burger not found", "id", id)

		return model.Burger{}, err
	}

	if err != nil {
		log.Info("fail to get burger", "id", id)

		return model.Burger{}, err
	}

	return burger, nil
}

func (s *Service) SaveBurger(ctx context.Context, burgerInfo serviceModel.BurgerInfo) error {
	const fn = "service.burger.save"

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	for _, ingredient := range burgerInfo.Ingredients {
		_, err := s.ingredientServer.GetIngredient(ctx, ingredient.IngredientId)
		if err != nil {
			log.Error("fail to get ingredient", sl.Err(err), "ingredient", ingredient)

			return err
		}
	}

	burgerInfo.Handle = strcase.ToKebab(burgerInfo.Title)
	burgerInfo.DataModified = time.Now()

	err := s.burgerRepository.SaveBurger(ctx, burgerInfo)
	if err != nil {
		log.Error("fail to save burger", sl.Err(err), "burgerInfo", burgerInfo)

		return err
	}

	return nil
}
