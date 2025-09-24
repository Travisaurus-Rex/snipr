package handlers

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
