package verify

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
)

type TokenService interface {
	VerifyEmail(ctx context.Context, token string) error
}

// New godoc
//
//	@Summary		verify user email
//	@Description	verify user email
//	@Tags			auth
//	@Accept			json
//
//	@Produce		json
//
//	@Param			token	path	string	true	"token"
//	@Success		204
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/auth/verify/{token} [get]
func New(tokenService TokenService) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.auth.validate.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		token := r.PathValue("token")
		if token == "" {
			log.Error("token is empty")
			return api.Error("token is empty", http.StatusBadRequest)
		}

		err := tokenService.VerifyEmail(ctx, token)
		if err != nil {
			if errors.Is(err, errs.ErrTokenNotFound) {
				log.Error("token not found", logger.Err(err))
				return api.Error(errs.ErrTokenNotFound.Error(), http.StatusNotFound)
			}
			log.Error("failed to validate token", logger.Err(err))
			return api.Error("failed to validate token", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}
