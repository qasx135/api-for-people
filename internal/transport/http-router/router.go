package http_router

import (
	userhandlers "api-for-people/internal/transport/http-router/user-handlers"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type Config struct {
	Address string `env:"ADDRESS,required" envDefault:"localhost:8080"`
}
type Router struct {
	Router *chi.Mux
	Config Config
}

func NewRouter(cfg Config, h *userhandlers.Handler) *Router {
	slog.Info("Loading router")
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Get("/persons", h.GetAll)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return &Router{Router: r, Config: cfg}
}

func (r *Router) Run(cfg Config, router *Router) {
	slog.Info("Starting server")
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router.Router,
	}
	fmt.Println("Listening on " + cfg.Address)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Error starting server", err)
	}
}
