package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
)

type SessionValidator interface {
	ValidateAdminSession(ctx context.Context, sessionId string) error
	ValidateUserSession(ctx context.Context, sessionId string) (string, error)
}

func Admin(sessionValidator SessionValidator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.Auth"
			ctx := r.Context()
			log := logger.FromCtx(ctx).With(slog.String("op", op))

			session, err := r.Cookie("session_id")
			if err != nil {
				log.Error("failed to get session", logger.Err(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			err = sessionValidator.ValidateAdminSession(ctx, session.Value)
			if err != nil {
				log.Error("failed to validate session", logger.Err(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func User(sessionValidator SessionValidator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.Auth"
			ctx := r.Context()
			log := logger.FromCtx(ctx).With(slog.String("op", op))

			session, err := r.Cookie("session_id")
			if err != nil {
				log.Error("failed to get session", logger.Err(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userId, err := sessionValidator.ValidateUserSession(ctx, session.Value)
			if err != nil {
				log.Error("failed to validate session", logger.Err(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "user_id", userId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
