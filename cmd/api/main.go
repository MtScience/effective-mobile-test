package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "subscriptions/internal/db/postgres"
	"subscriptions/internal/logging"
	repo "subscriptions/internal/repository/postgres"
	transporthttp "subscriptions/internal/transport/http"
	swaggerdocs "subscriptions/internal/transport/http/docs"
	"subscriptions/internal/transport/http/handlers"
	"subscriptions/internal/usecase/subscription"
)

func main() {

	cfg := loadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.Open(ctx, cfg.DatabaseDSN)
	if err != nil {
		logging.Logger.Error("open database connection", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	subRepo := repo.New(pool)
	subUsecase := subscription.NewService(subRepo)
	subHandler := handlers.NewSubscriptionHandler(subUsecase)
	docsHandler := swaggerdocs.NewHandler()
	router := transporthttp.NewRouter(subHandler, docsHandler)

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logging.Logger.Error("shutdown http server", "error", err)
		}
	}()

	logging.Logger.Info("http server started", "addr", cfg.HTTPAddress)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logging.Logger.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}

type config struct {
	HTTPAddress string
	DatabaseDSN string
}

func loadConfig() config {
	httpAddress := os.Getenv("HTTP_ADDRESS")
	if httpAddress == "" {
		logging.Logger.Warn("environment variable 'HTTP_ADDRESS' is not set. Using default value")
		httpAddress = ":8080"
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	for _, part := range []string{dbUser, dbPass, dbName, dbPort} {
		if part == "" {
			panic(fmt.Errorf("invalid database DSN"))
		}
	}

	dbDsn := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable", dbUser, dbPass, dbPort, dbName)

	return config{
		HTTPAddress: httpAddress,
		DatabaseDSN: dbDsn,
	}
}
