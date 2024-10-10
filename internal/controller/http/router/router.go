package router

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"practice/internal/controller/http/handler"
	"practice/internal/pkg/config"
	"strings"

	swagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

type Options struct {
	fx.In
	fx.Lifecycle
	Config  *config.Config
	Logger  *slog.Logger
	Handler *handler.Handler
}

var Module = fx.Options(
	fx.Invoke(New),
)

func New(opts Options) {
	router := chi.NewRouter()

	router.Mount("/docs", swagger.WrapHandler)

	router.Route("/user", func(r chi.Router) {
		r.Get("/{id}", opts.Handler.GetUser)
		r.Post("/", opts.Handler.CreateUser)
		r.Put("/{id}", opts.Handler.UpdateUser)
		r.Delete("/{id}", opts.Handler.DeleteUser)
	})

	router.Route("/computer", func(r chi.Router) {
		r.Get("/{id}", opts.Handler.GetComputer)
		r.Post("/", opts.Handler.CreateComputer)
		r.Put("/{id}", opts.Handler.UpdateComputer)
		r.Delete("/{id}", opts.Handler.DeleteComputer)
		r.Get("/", opts.Handler.ListComputers)
	})

	server := http.Server{
		Addr:         fmt.Sprint(opts.Config.ADDRESS),
		Handler:      router,
		ReadTimeout:  opts.Config.ReadTimeout,
		WriteTimeout: opts.Config.WriteTimeout,
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: onStart(&server, opts.Config, opts.Logger),
		OnStop:  onStop(&server, opts.Logger),
	})
}

func onStart(srv *http.Server, cfg *config.Config, log *slog.Logger) func(_ context.Context) error {
	return func(_ context.Context) error {
		log.Info("starting server on %s", cfg.ADDRESS)
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				panic("failed to start server: " + err.Error())
			}
		}()
		return nil
	}
}

func onStop(srv *http.Server, log *slog.Logger) func(_ context.Context) error {
	return func(ctx context.Context) error {
		log.Info("shutdown server by signal")
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server forced to shutdown: %v", err)
		}
		return nil
	}
}

func displayRoutes(r *chi.Mux, l *slog.Logger) {
	routes := r.Routes()
	for _, route := range routes {
		if strings.Contains(route.Pattern, "/docs") {
			continue
		}

		for handle := range route.Handlers {
			l.Debug("http controller. mapped [%v] %v route", handle, route.Pattern)
		}
	}
}
