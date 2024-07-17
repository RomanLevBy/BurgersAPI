package all

import (
	"context"
	"fmt"
	resp "github.com/RomanLevBy/BurgersAPI/internal/lib/api/response"
	"github.com/RomanLevBy/BurgersAPI/internal/lib/logger/sl"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/ingredient/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Ingredients []model.Ingredient `json:"ingredient,required"`
}

type IngredientsProvider interface {
	GetAllIngredients(ctx context.Context, param serviceModel.FetchParam) ([]model.Ingredient, uint64, error)
}

func New(log *slog.Logger, burgersProvider IngredientsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.ingredient.all.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		title := r.URL.Query().Get("i")
		if title == "" {
			log.Error("i is a required field")

			render.JSON(w, r, resp.Error("i is a required field."))

			return
		}

		params := serviceModel.FetchParam{
			Title: title,
		}

		ingredients, nextCursor, err := burgersProvider.GetAllIngredients(r.Context(), params)
		if err != nil {
			log.Info("fail to get ingredient", sl.Err(err))

			render.JSON(w, r, resp.Error("Fail to get ingredient."))

			return
		}

		log.Info("ingredient found", "ingredient", ingredients)

		w.Header().Add("X-NextCursor", fmt.Sprintf("%d", nextCursor))

		render.JSON(w, r, Response{
			Response:    resp.OK(),
			Ingredients: ingredients,
		})
	}
}
