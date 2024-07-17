package all

import (
	"context"
	"fmt"
	resp "github.com/RomanLevBy/BurgersAPI/internal/lib/api/response"
	"github.com/RomanLevBy/BurgersAPI/internal/lib/logger/sl"
	"github.com/RomanLevBy/BurgersAPI/internal/model"
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/burger/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

type Response struct {
	resp.Response
	Burgers []model.Burger `json:"burgers,required"`
}

type BurgersProvider interface {
	GetAllBurgers(ctx context.Context, param serviceModel.FetchParam) ([]model.Burger, uint64, error)
}

func New(log *slog.Logger, burgersProvider BurgersProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.burger.all.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil && limitStr != "" {
			log.Error("incorrect limit parameter")

			render.JSON(w, r, resp.Error("Incorrect limit parameter."))

			return
		}

		if limit > maxLimit {
			log.Error("limit is more than max limit", slog.Int("limit", limit))

			render.JSON(w, r, resp.Error(fmt.Sprintf("Limit is more than max limit: %d", maxLimit)))

			return
		}

		if limit == 0 {
			limit = defaultLimit
		}

		cursorStr := r.URL.Query().Get("cursor")
		cursor, err := strconv.Atoi(cursorStr)
		if err != nil && cursorStr != "" {
			log.Error("cursor parameter is invalid")

			render.JSON(w, r, resp.Error("Incorrect cursor parameter."))

			return
		}

		title := r.URL.Query().Get("s")
		titlePath := r.URL.Query().Get("f")

		params := serviceModel.FetchParam{
			Title:     title,
			TitlePath: titlePath,
			Limit:     uint64(limit),
			CursorID:  uint64(cursor),
		}

		burgers, nextCursor, err := burgersProvider.GetAllBurgers(r.Context(), params)
		if err != nil {
			log.Info("fail to get burgers", sl.Err(err))

			render.JSON(w, r, resp.Error("Fail to get burgers."))

			return
		}

		log.Info("burgers found", "burgers", burgers)

		w.Header().Add("X-NextCursor", fmt.Sprintf("%d", nextCursor))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Burgers:  burgers,
		})
	}
}
