package handlers

import (
	"net/http"

	"github.com/Travisaurus-Rex/snipr/internal/db"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, url.LongURL, http.StatusFound)
}
