package get

import (
	"context"
	"errors"
	iErr "github.com/RomanLevBy/BurgersAPI/internal/error"
	resp "github.com/RomanLevBy/BurgersAPI/internal/lib/api/response"
	"github.com/RomanLevBy/BurgersAPI/internal/lib/logger/sl"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	resp.Response
	Burger model.Burger `json:"burger,required"`
}

type BurgerProvider interface {
	GetBurger(ctx context.Context, id int) (model.Burger, error)
}

func New(log *slog.Logger, burgerProvider BurgerProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.burger.get.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			log.Error("Id is empty")

			render.JSON(w, r, resp.Error("Id is empty."))

			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("id is not valid")

			render.JSON(w, r, resp.Error("id is not valid."))

			return
		}

		burger, err := burgerProvider.GetBurger(r.Context(), id)
		if errors.Is(err, iErr.ErrBurgerNotFound) {
			log.Info("burger not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("Not found."))

			return
		}

		if err != nil {
			log.Info("fail to get burger", slog.String("id", idParam), sl.Err(err))

			render.JSON(w, r, resp.Error("Fail to get burger."))

			return
		}

		log.Info("burger found", slog.String("burger", idParam), "burger", burger)

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Burger:   burger,
		})
	}
}
