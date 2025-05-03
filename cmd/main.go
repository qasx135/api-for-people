package main

import (
	"api-for-people/internal/config"
	http_router "api-for-people/internal/transport/http-router"
	user_handlers "api-for-people/internal/transport/http-router/user-handlers"
	"api-for-people/internal/user/repository"
	"api-for-people/internal/user/service"
	"api-for-people/pkg/postgres"
	"context"
	"fmt"
	"log/slog"
)

func main() {
	cfg := config.NewConfig()
	fmt.Println(cfg)
	ctx := context.Background()
	db, _ := postgres.New(ctx, &cfg.PostgresConfig)
	repo := repository.NewRepository(db)
	newService := service.NewService(repo)
	handler := user_handlers.NewHandler(ctx, newService)
	router := http_router.NewRouter(cfg.RouterConfig, handler)
	slog.Info("Starting server")
	router.Run(cfg.RouterConfig, router)
	defer slog.Info("Stopping application")
}
