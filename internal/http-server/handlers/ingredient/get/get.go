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
	Ingredient model.Ingredient `json:"ingredient,required"`
}

type IngredientProvider interface {
	GetIngredient(ctx context.Context, id int) (model.Ingredient, error)
}

func New(log *slog.Logger, burgerProvider IngredientProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.ingredient.get.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			log.Error("id is empty")

			render.JSON(w, r, resp.Error("Id is empty."))

			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("id is not valid")

			render.JSON(w, r, resp.Error("Id is not valid."))

			return
		}

		ingredient, err := burgerProvider.GetIngredient(r.Context(), id)
		if errors.Is(err, iErr.ErrIngredientNotFound) {
			log.Info("ingredient not found", slog.Int("id", id))

			render.JSON(w, r, resp.Error("Not found."))

			return
		}

		if err != nil {
			log.Info("fail to get ingredient", slog.String("id", idParam), sl.Err(err))

			render.JSON(w, r, resp.Error("Fail to get ingredient."))

			return
		}

		log.Info("ingredient found", slog.String("ingredient", idParam), "ingredient", ingredient)

		render.JSON(w, r, Response{
			Response:   resp.OK(),
			Ingredient: ingredient,
		})
	}
}
