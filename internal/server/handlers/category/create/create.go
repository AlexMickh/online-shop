package create_category

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name string `json:"name" validate:"required,min=3"`
}

type Response struct {
	ID string `json:"id"`
}

type CategoryCreator interface {
	CreateCategory(ctx context.Context, name string) (string, error)
}

// New godoc
//
//	@Summary		create new category
//	@Description	create new category
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			name	body		string	true	"category name"
//	@Success		201		{object}	Response
//	@Failure		400		{object}	api.ErrorResponse
//	@Failure		500		{object}	api.ErrorResponse
//	@Security		SessionAuth
//	@Router			/admin/create-category [post]
func New(categoryCreator CategoryCreator, validator *validator.Validate) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.category.create.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode body", logger.Err(err))
			return api.Error("failed to decode body", http.StatusBadRequest)
		}
		defer r.Body.Close()

		if err := validator.Struct(&req); err != nil {
			log.Error("failed to validate request", logger.Err(err))
			return api.Error("failed to validate request", http.StatusBadRequest)
		}

		id, err := categoryCreator.CreateCategory(ctx, req.Name)
		if err != nil && !errors.Is(err, errs.ErrFailedToCash) {
			if errors.Is(err, errs.ErrCategoryAlreadyExists) {
				log.Error("category already exists", logger.Err(err))
				return api.Error(errs.ErrCategoryAlreadyExists.Error(), http.StatusBadRequest)
			}
			log.Error("failed to create category", logger.Err(err))
			return api.Error("failed to create category", http.StatusInternalServerError)
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			ID: id,
		})

		return nil
	}
}
