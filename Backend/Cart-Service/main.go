package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kidkon/ecommerce/common/response"
)

const serviceName = "cart-service"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, "ok", map[string]string{"service": serviceName})
	})

	port := getenv("PORT", "8083")
	log.Printf("[%s] listening on :%s", serviceName, port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
