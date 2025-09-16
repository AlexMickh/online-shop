package get_category

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
)

type CategoryProvider interface {
	AllCategories(ctx context.Context) ([]models.Category, error)
}

type category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Categories []category `json:"categories"`
}

// New godoc
//
//	@Summary		returns all categories
//	@Description	returns all categories
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/category [get]
func New(categoryProvider CategoryProvider) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.category.get.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		categoriesInfo, err := categoryProvider.AllCategories(ctx)
		if err != nil {
			log.Error("failed to get all categories", logger.Err(err))
			return api.Error("failed to get all categories", http.StatusInternalServerError)
		}

		categories := make([]category, 0, len(categoriesInfo))
		for _, c := range categoriesInfo {
			category := category{
				ID:   c.ID,
				Name: c.Name,
			}
			categories = append(categories, category)
		}

		render.JSON(w, r, Response{
			Categories: categories,
		})

		return nil
	}
}
