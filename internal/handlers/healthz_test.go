package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	Healthz(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Result().StatusCode)
	}

	body := w.Body.String()
	expectedBody := "ok"

	if body != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}

	t.Logf("Request method: %s", req.Method)
	t.Logf("Request URL path: %s", req.URL.Path)
	t.Logf("Request headers: %v", req.Header)
}
