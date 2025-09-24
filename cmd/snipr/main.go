package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Travisaurus-Rex/snipr/internal/db"
	"github.com/Travisaurus-Rex/snipr/internal/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	db.Connect()
	defer db.Pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handlers.Shorten)
	mux.HandleFunc("/healthz", handlers.Healthz)
	mux.HandleFunc("/", handlers.Redirect)

	port := getenv("APP_PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
