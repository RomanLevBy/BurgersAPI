package burger

import (
	"context"
	"errors"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	"github.com/RomanLevBy/BurgersAPI/internal/lib/logger/sl"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iancoleman/strcase"
	"log/slog"
	"time"
)

type Repository interface {
	GetBurger(ctx context.Context, id int) (model.Burger, error)
	SaveBurger(ctx context.Context, burgerInfo model.BurgerInfo) error
}

type Service struct {
	logger           *slog.Logger
	burgerRepository Repository
}

func New(repo Repository, logger *slog.Logger) *Service {
	return &Service{burgerRepository: repo, logger: logger}
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

func (s *Service) SaveBurger(ctx context.Context, burgerInfo model.BurgerInfo) error {
	const fn = "service.burger.save"

	log := s.logger.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	burgerInfo.Handle = strcase.ToKebab(burgerInfo.Title)
	burgerInfo.DataModified = time.Now()

	err := s.burgerRepository.SaveBurger(ctx, burgerInfo)
	if err != nil {
		log.Error("fail to save burger", sl.Err(err), "burgerInfo", burgerInfo)

		return err
	}

	return nil
}
