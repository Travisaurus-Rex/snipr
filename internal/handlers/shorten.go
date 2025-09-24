package handlers

import (
	"fmt"
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
	fmt.Fprintf(w, `{"id":%d,"short_code":"%s","long_url":"%s"}`, url.ID, url.ShortCode, url.LongURL)
}
