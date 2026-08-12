package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleUsageRecords(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	limit := 100
	offset := 0

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	records, err := s.store.ListUsageRecords(r.Context(), id, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, records)
}

func (s *Server) handleUsageDailySummary(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	summaries, err := s.store.ListUsageDailySummary(r.Context(), id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, summaries)
}
