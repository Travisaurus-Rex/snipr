package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Travisaurus-Rex/snipr/internal/db"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load()
	db.Connect()
	code := m.Run()

	if db.Pool != nil {
		db.Pool.Close()
	}

	os.Exit(code)
}

func TestShorten_DBUnavailable(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}
	db.Pool = nil

	r := httptest.NewRequest(http.MethodGet, "/shorten?long_url=https://example.com", nil)
	w := httptest.NewRecorder()

	Shorten(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	expected := "error while shortening the provided url\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

func TestShorten_MissingLongURL(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	w := httptest.NewRecorder()

	Shorten(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expected := "long_url missing in query params\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

func TestShorten_Success(t *testing.T) {
	db.Connect()

	r := httptest.NewRequest(http.MethodGet, "/shorten?long_url=https://www.fakeurl.com", nil)
	w := httptest.NewRecorder()

	Shorten(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected content type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if res["short_code"] == "" {
		t.Error("expected short_code to not be empty")
	}

	if res["long_url"] != "https://www.fakeurl.com" {
		t.Errorf("expected long_url to be %q, got %q", "https://www.fakeurl.com", res["long_url"])
	}
}
