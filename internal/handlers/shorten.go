package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Travisaurus-Rex/snipr/internal/db"
	"github.com/Travisaurus-Rex/snipr/internal/shortener"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
	longURL := r.URL.Query().Get("long_url")
	if longURL == "" {
		http.Error(w, "long_url missing in query params", http.StatusBadRequest)
		return
	}

	shortCode := shortener.GenerateCode()

	url, err := db.InsertURL(shortCode, longURL)
	if err != nil {
		log.Printf("db error: %v", err)
		http.Error(w, "error while shortening the provided url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := map[string]interface{}{
		"id":         url.ID,
		"short_code": url.ShortCode,
		"long_url":   url.LongURL,
		"created_at": url.CreatedAt,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("json encoding error: %v", err)
	}
}
