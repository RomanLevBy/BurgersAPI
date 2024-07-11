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
	Category model.Category `json:"category,required"`
}

type CategoryProvider interface {
	GetCategory(ctx context.Context, id int) (model.Category, error)
}

func New(log *slog.Logger, categorySaver CategoryProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.category.get.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			log.Error("Id is empty")

			render.JSON(w, r, resp.Error("Id is empty"))

			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("Id is not valid")

			render.JSON(w, r, resp.Error("Id is not valid"))

			return
		}

		category, err := categorySaver.GetCategory(context.Background(), id)
		if errors.Is(err, iErr.ErrCategoryNotFound) {
			log.Info("category not found", "id", id)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Info("fail to get category", slog.String("id", idParam), sl.Err(err))

			render.JSON(w, r, resp.Error("Fail to get category"))

			return
		}

		log.Info("category found", slog.String("category", idParam))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Category: category,
		})
	}
}
