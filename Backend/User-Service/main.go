package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/kidkon/ecommerce/user-service/internal/authmw"
	"github.com/kidkon/ecommerce/common/config"
	"github.com/kidkon/ecommerce/common/database"
	"github.com/kidkon/ecommerce/common/logger"
	"github.com/kidkon/ecommerce/common/response"
	"github.com/kidkon/ecommerce/user-service/internal/handler"
	"github.com/kidkon/ecommerce/user-service/internal/repository"
	"github.com/kidkon/ecommerce/user-service/internal/service"
)

func main() {
	config.LoadDotenv()
	logger.Setup("user-service", config.Get("LOG_LEVEL", "info"))

	ctx := context.Background()
	pool, err := database.ConnectFromEnv(ctx)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected")

	jwtSecret := config.Get("JWT_SECRET", "dev-secret")

	// wire: repository -> service -> handler
	repo := repository.NewUserRepository(pool)
	authSvc := service.NewAuthService(repo, jwtSecret)
	authH := handler.NewAuthHandler(authSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"service": "user-service"})
	})

	// public routes (ไม่ต้อง token)
	mux.HandleFunc("POST /api/auth/register", authH.Register)
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// protected route: ครอบด้วย RequireAuth → /me แสดงข้อมูลจาก token
	mux.Handle("GET /api/auth/me", authmw.RequireAuth(jwtSecret)(http.HandlerFunc(authH.Me)))
	mux.Handle("GET /api/auth/me/profile", authmw.RequireAuth(jwtSecret)(http.HandlerFunc(authH.Profile)))

	port := config.Get("PORT", "8081")
	slog.Info("service starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server stopped", "err", err)
	}
}
