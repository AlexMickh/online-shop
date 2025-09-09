package register

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/internal/consts"
	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Login    string `json:"login" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type Response struct {
	ID string `json:"id"`
}

type Registerer interface {
	Register(ctx context.Context, login, email, password string) (string, error)
}

type TokenCreator interface {
	CreateToken(ctx context.Context, userId, tokenType string) (string, error)
}

type VerificationSender interface {
	SendVerification(to string, token, login string) error
}

func New(
	validator *validator.Validate,
	registerer Registerer,
	tokenCreator TokenCreator,
	verificationSender VerificationSender,
) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.auth.register.New"
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

		id, err := registerer.Register(ctx, req.Login, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, errs.ErrUserAlreadyExists) {
				log.Error("user already exists", logger.Err(err))
				return api.Error(errs.ErrUserAlreadyExists.Error(), http.StatusBadRequest)
			}
			log.Error("failed to register user", logger.Err(err))
			return api.Error("failed to register user", http.StatusInternalServerError)
		}

		token, err := tokenCreator.CreateToken(ctx, id, consts.TokenTypeEmailVerify)
		if err != nil {
			log.Error("failed to create token", logger.Err(err))
			return api.Error("failed to create token", http.StatusInternalServerError)
		}

		err = verificationSender.SendVerification(req.Email, token, req.Login)
		if err != nil {
			log.Error("failed to send email", logger.Err(err))
			return api.Error("failed to send email", http.StatusInternalServerError)
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			ID: id,
		})

		return nil
	}
}
