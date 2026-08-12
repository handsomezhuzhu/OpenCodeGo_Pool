package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchPage_CookieExpiredWhenNoFlightMarker(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>please log in</body></html>"))
	}))
	defer srv.Close()

	c := NewClient(5 * time.Second)
	c.baseURL = srv.URL

	_, err := c.fetchPage("session=abc", "/workspace/ws1/go")
	if err == nil {
		t.Fatal("expected error for page without $R marker")
	}
	if !strings.Contains(err.Error(), "cookie expired or invalid") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFetchPage_PassesThroughWhenFlightMarkerPresent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><script>$R[0] = {}</script></body></html>`))
	}))
	defer srv.Close()

	c := NewClient(5 * time.Second)
	c.baseURL = srv.URL

	body, err := c.fetchPage("session=abc", "/workspace/ws1/go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, "$R") {
		t.Error("expected body to be returned unchanged")
	}
}

func TestFetchPage_RedirectTriggersAuthError(t *testing.T) {
	cases := []struct {
		name   string
		status int
	}{
		{"moved permanently", http.StatusMovedPermanently},
		{"found", http.StatusFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/login", tc.status)
			}))
			defer srv.Close()

			c := NewClient(5 * time.Second)
			c.baseURL = srv.URL

			_, err := c.fetchPage("session=abc", "/workspace/ws1/go")
			if err == nil {
				t.Fatal("expected error for redirect response")
			}
			if !strings.Contains(err.Error(), "cookie expired (redirect to login)") {
				t.Errorf("unexpected error message: %v", err)
			}
		})
	}
}

func TestFetchPage_ExistingAuthStatusBehaviorUnchanged(t *testing.T) {
	cases := []struct {
		name   string
		status int
	}{
		{"unauthorized", http.StatusUnauthorized},
		{"forbidden", http.StatusForbidden},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()

			c := NewClient(5 * time.Second)
			c.baseURL = srv.URL

			_, err := c.fetchPage("session=abc", "/workspace/ws1/go")
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), "cookie expired") {
				t.Errorf("unexpected error message: %v", err)
			}
		})
	}
}
