package login

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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type Loginer interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type SessionCreator interface {
	Create(w http.ResponseWriter, sessionId string)
}

// New godoc
//
//	@Summary		login user
//	@Description	login user
//	@Tags			auth
//	@Accept			json
//
//	@Produce		json
//
//	@Param			email		body	string	true	"User email"	Format(email)
//	@Param			password	body	string	true	"User password"
//	@Success		201
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		404	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/auth/login [post]
func New(loginer Loginer, validator validator.Validate, sessionCreator SessionCreator) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.auth.login.New"
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

		sessionId, err := loginer.Login(ctx, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				log.Error("user not found", logger.Err(err))
				return api.Error(errs.ErrUserNotFound.Error(), http.StatusNotFound)
			}

			log.Error("failed to login user", logger.Err(err))
			return api.Error("failed to login user", http.StatusInternalServerError)
		}

		sessionCreator.Create(w, sessionId)

		return nil
	}
}
