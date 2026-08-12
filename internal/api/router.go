package api

import (
	"io/fs"
	"net/http"
	"strings"

	"OpenCodeGoPool/internal/cpa"
	"OpenCodeGoPool/internal/scheduler"
	"OpenCodeGoPool/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	store     store.Store
	scheduler *scheduler.Scheduler
	cpa       *cpa.Syncer
	password  string
	sessions  *SessionStore
}

func NewServer(s store.Store, sched *scheduler.Scheduler, cpaSyncer *cpa.Syncer, password string) *Server {
	return &Server{
		store:     s,
		scheduler: sched,
		cpa:       cpaSyncer,
		password:  password,
		sessions:  NewSessionStore(),
	}
}

func (s *Server) Router(frontendFS fs.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", s.handleLogin)
		r.Get("/auth/check", s.handleAuthCheck)

		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware)

			r.Get("/dashboard", s.handleDashboard)

			r.Get("/accounts", s.handleListAccounts)
			r.Post("/accounts", s.handleCreateAccount)
			r.Get("/accounts/{id}", s.handleGetAccount)
			r.Put("/accounts/{id}", s.handleUpdateAccount)
			r.Delete("/accounts/{id}", s.handleDeleteAccount)
			r.Patch("/accounts/{id}/status", s.handleToggleStatus)
			r.Get("/accounts/{id}/quota", s.handleQuotaHistory)
			r.Get("/accounts/{id}/usage/daily", s.handleUsageDailySummary)
			r.Get("/accounts/{id}/usage", s.handleUsageRecords)
			r.Post("/accounts/{id}/refresh", s.handleRefreshOne)

			r.Post("/refresh", s.handleRefreshAll)

			r.Post("/sync/cpa", s.handleSyncCPA)
			r.Get("/sync/cpa/status", s.handleSyncStatus)

			r.Get("/settings/cpa", s.handleGetCPASettings)
			r.Put("/settings/cpa", s.handleUpdateCPASettings)
		})
	})

	if frontendFS != nil {
		spaHandler := spaFileServer(frontendFS)
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			spaHandler.ServeHTTP(w, r)
		})
	}

	return r
}

func spaFileServer(frontendFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(frontendFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		_, err := fs.Stat(frontendFS, path)
		if err != nil {
			r.URL.Path = "/"
			path = "index.html"
		}

		if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}

		fileServer.ServeHTTP(w, r)
	})
}
