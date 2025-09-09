package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexMickh/coledzh-shop-backend/internal/app"
	"github.com/AlexMickh/coledzh-shop-backend/internal/config"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	file, err := os.OpenFile(cfg.Env+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log := logger.New(cfg.Env, io.MultiWriter(os.Stdout, file))

	log.Info("logger is working", slog.String("env", cfg.Env))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = logger.ContextWithLogger(ctx, log)

	app := app.New(ctx, cfg)
	app.Run(ctx)
	defer app.GracefulStop(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	close(stop)
	logger.FromCtx(ctx).Info("server stopped")
}
