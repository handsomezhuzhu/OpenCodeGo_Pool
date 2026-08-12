package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleSyncCPA(w http.ResponseWriter, r *http.Request) {
	if err := s.cpa.Sync(context.Background()); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"status": "synced"})
}

func (s *Server) handleSyncStatus(w http.ResponseWriter, r *http.Request) {
	log, err := s.store.GetLatestSyncLog(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, log)
}

func (s *Server) handleRefreshAll(w http.ResponseWriter, r *http.Request) {
	go s.scheduler.ScrapeAll(context.Background())
	jsonOK(w, map[string]string{"status": "refreshing"})
}

func (s *Server) handleRefreshOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	go s.scheduler.ScrapeOne(context.Background(), id)
	jsonOK(w, map[string]string{"status": "refreshing"})
}
