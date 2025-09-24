package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Travisaurus-Rex/snipr/internal/db"
	"github.com/Travisaurus-Rex/snipr/internal/shortener"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	db.Connect()
	defer db.Pool.Close()

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		longURL := r.URL.Query().Get("long_url")
		if longURL == "" {
			http.Error(w, "long_url missing in query params", http.StatusBadRequest)
			return
		}

		shortCode := shortener.GenerateCode()

		url, err := db.InsertURL(shortCode, longURL)
		if err != nil {
			http.Error(w, "error while shortening the provided url", http.StatusInternalServerError)
			log.Fatalf("db error: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"id":%d,"short_code":"%s","long_url":"%s"}`, url.ID, url.ShortCode, url.LongURL)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.URL.Path[1:]
		if shortCode == "" {
			http.Error(w, "short code not provided", http.StatusBadRequest)
			return
		}

		url, err := db.GetURLByCode(shortCode)
		if err != nil {
			if err.Error() == "no URL found for that short code" {
				http.NotFound(w, r)
				return
			}

			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url.LongURL, http.StatusOK)
	})

	port := getenv("APP_PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
