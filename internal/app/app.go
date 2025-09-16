package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/AlexMickh/coledzh-shop-backend/internal/config"
	category_repository "github.com/AlexMickh/coledzh-shop-backend/internal/repository/postgres/category"
	token_repository "github.com/AlexMickh/coledzh-shop-backend/internal/repository/postgres/token"
	user_repository "github.com/AlexMickh/coledzh-shop-backend/internal/repository/postgres/user"
	category_cash "github.com/AlexMickh/coledzh-shop-backend/internal/repository/redis/category"
	session_cash "github.com/AlexMickh/coledzh-shop-backend/internal/repository/redis/session"
	"github.com/AlexMickh/coledzh-shop-backend/internal/server"
	auth_service "github.com/AlexMickh/coledzh-shop-backend/internal/services/auth"
	category_service "github.com/AlexMickh/coledzh-shop-backend/internal/services/category"
	token_service "github.com/AlexMickh/coledzh-shop-backend/internal/services/token"
	user_service "github.com/AlexMickh/coledzh-shop-backend/internal/services/user"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/clients/postgresql"
	redis_client "github.com/AlexMickh/coledzh-shop-backend/pkg/clients/redis"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	srv *server.Server
	db  *pgxpool.Pool
	rdb *redis.Client
}

func New(ctx context.Context, cfg *config.Config) *App {
	const op = "app.New"

	log := logger.FromCtx(ctx).With(slog.String("op", op))

	log.Info("initing postgres")
	db, err := postgresql.New(
		ctx,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.MinPools,
		cfg.DB.MaxPools,
		cfg.DB.MigrationsPath,
	)
	if err != nil {
		log.Error("failed to init postgres", logger.Err(err))
		os.Exit(1)
	}
	userRepository := user_repository.New(db)
	tokenRepository := token_repository.New(db)
	categoryRepository := category_repository.New(db)

	log.Info("initing redis")
	cash, err := redis_client.New(
		ctx,
		fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.User,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Error("failed to init redis", logger.Err(err))
		os.Exit(1)
	}
	sessionCash := session_cash.New(cash, cfg.Redis.Expiration)
	categoryCash := category_cash.New(cash, cfg.Redis.Expiration)

	log.Info("initing service layer")
	authService := auth_service.New(userRepository, sessionCash)
	tokenService := token_service.New(tokenRepository, authService)
	categoryService := category_service.New(categoryRepository, categoryCash)
	userService := user_service.New(sessionCash)

	log.Info("initing server")
	srv := server.New(
		ctx,
		cfg.Server,
		authService,
		cfg.Mail,
		tokenService,
		categoryService,
		userService,
	)

	return &App{
		srv: srv,
		db:  db,
		rdb: cash,
	}
}

func (a *App) Run(ctx context.Context) {
	const op = "app.Run"

	log := logger.FromCtx(ctx).With(slog.String("op", op))

	go func() {
		if err := a.srv.Run(ctx); err != nil {
			log.Error("failed to start server", logger.Err(err))
			os.Exit(1)
		}
	}()

	log.Info("server started", slog.String("addr", a.srv.Addr()))
}

func (a *App) GracefulStop(ctx context.Context) {
	a.srv.GracefulStop(ctx)
	a.db.Close()
	a.rdb.Close()
}
