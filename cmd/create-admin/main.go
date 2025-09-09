package main

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/config"
	user_repository "github.com/AlexMickh/coledzh-shop-backend/internal/repository/postgres/user"
	session_cash "github.com/AlexMickh/coledzh-shop-backend/internal/repository/redis/session"
	auth_service "github.com/AlexMickh/coledzh-shop-backend/internal/services/auth"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/clients/postgresql"
	redis_client "github.com/AlexMickh/coledzh-shop-backend/pkg/clients/redis"
)

func main() {
	var (
		login    string
		email    string
		password string
	)

	fmt.Print("Input login: ")
	fmt.Scan(&login)
	fmt.Print("Input email: ")
	fmt.Scan(&email)
	fmt.Print("Input password: ")
	fmt.Scan(&password)

	cfg := config.MustLoad()

	db, err := postgresql.New(
		context.Background(),
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
		panic(err)
	}
	defer db.Close()
	userRepository := user_repository.New(db)

	cash, err := redis_client.New(
		context.Background(),
		fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.User,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		panic(err)
	}
	sessionCash := session_cash.New(cash, cfg.Redis.Expiration)

	authService := auth_service.New(userRepository, sessionCash)

	err = authService.RegisterAdmin(context.Background(), login, email, password)
	if err != nil {
		panic(err)
	}

	fmt.Printf("admin with login %s successfully created", login)
}
