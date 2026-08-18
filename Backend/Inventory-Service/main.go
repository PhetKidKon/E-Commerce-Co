package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/kidkon/ecommerce/common/config"
	"github.com/kidkon/ecommerce/common/database"
	"github.com/kidkon/ecommerce/common/logger"
	"github.com/kidkon/ecommerce/common/response"
)

func main() {
	// load Backend/.env (local dev) then set up the shared logger
	config.LoadDotenv()
	logger.Setup("inventory-service", config.Get("LOG_LEVEL", "info"))

	// connect to Postgres at startup (no handlers yet — just prove connectivity)
	ctx := context.Background()
	pool, err := database.ConnectFromEnv(ctx)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1) // fail fast: a service without its DB can't do its job
	}
	defer pool.Close()
	slog.Info("database connected")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"service": "inventory-service"})
	})
	mux.HandleFunc("/api/inventory", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"service": "inventory-service"})
	})

	port := config.Get("PORT", "8084")
	slog.Info("service starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server stopped", "err", err)
	}
}
