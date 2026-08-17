package main

import (
	"log/slog"
	"net/http"

	"github.com/kidkon/ecommerce/common/config"
	"github.com/kidkon/ecommerce/common/logger"
	"github.com/kidkon/ecommerce/common/response"
)

func main() {
	// wire shared common packages (bare — just proves the connection works)
	logger.Setup("cart-service", config.Get("LOG_LEVEL", "info"))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"service": "cart-service"})
	})

	port := config.Get("PORT", "8083")
	slog.Info("service starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server stopped", "err", err)
	}
}
