package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorishak/url-shortener/internal/config"
	"github.com/gorishak/url-shortener/internal/http-server/handlers/redirect"
	"github.com/gorishak/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/gorishak/url-shortener/internal/http-server/middleware/logger"
	"github.com/gorishak/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/gorishak/url-shortener/internal/lib/logger/sl"
	"github.com/gorishak/url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)
	fmt.Println("config loaded")

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", "env", cfg.Env)
	log.Debug("debug messages are enabled")

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, "chi render"
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // это для id request
	router.Use(middleware.RealIP)    // это для ip address
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // это для recovery
	router.Use(middleware.URLFormat) // это для url format

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{cfg.User: cfg.Password}))
		r.Post("/", save.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Info("server stopped")
	// router.Use(middleware.Logger)    // это для логирования

	// TODO: run server:
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
