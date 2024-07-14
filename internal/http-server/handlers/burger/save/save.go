package save

import (
	"context"
	"errors"
	"github.com/RomanLevBy/BurgersAPI/internal/converter"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	resp "github.com/RomanLevBy/BurgersAPI/internal/lib/api/response"
	"github.com/RomanLevBy/BurgersAPI/internal/lib/logger/sl"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
}

type BurgerSaver interface {
	SaveBurger(ctx context.Context, burgerInfo model.BurgerInfo) error
}

func New(log *slog.Logger, burgerSaver BurgerSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.burger.save.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req model.BurgerRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to decode request."))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("Invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		burgerInfo := converter.ToBurgerInfoFromRequest(req)
		err = burgerSaver.SaveBurger(r.Context(), burgerInfo)
		if errors.Is(err, iErr.ErrBurgerExists) {
			log.Info("burger already exists", "burgerInfo", burgerInfo)

			render.JSON(w, r, resp.Error("Burger already exists."))

			return
		}

		if err != nil {
			log.Error("failed to add burger", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to add burger."))

			return
		}

		log.Info("burger added")

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
