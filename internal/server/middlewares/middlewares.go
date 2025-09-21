package middlewares

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strings"

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

func IPFilterMiddleware(allowedCIDRs []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remoteIP := r.RemoteAddr
			log.Printf("Initial remote IP: %s", remoteIP)

			// Проверяем X-Real-IP, если доступен.
			if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				log.Printf("Using X-Real-IP header: %s", realIP)
				remoteIP = realIP
			}

			// Разделяем адрес на хост и порт.
			var host string
			if strings.Contains(remoteIP, ":") {
				var err error
				host, _, err = net.SplitHostPort(remoteIP)
				if err != nil {
					http.Error(w, "Invalid remote IP address", http.StatusBadRequest)
					return
				}
			} else {
				host = remoteIP
			}

			// Проверяем, разрешен ли IP-адрес.
			if !isIPAllowed(host, allowedCIDRs) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Передаем управление дальше.
			next.ServeHTTP(w, r)
		})
	}
}

// Проверяет, входит ли IP-адрес в разрешенные диапазоны.
func isIPAllowed(ip string, allowedCIDRs []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range allowedCIDRs {
		_, allowedNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if allowedNet.Contains(parsedIP) {
			return true
		}
	}
	return false
}
