package api

import (
	"encoding/json"
	"net/http"

	"OpenCodeGoPool/internal/model"
)

func (s *Server) handleGetCPASettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.store.GetCPASettings(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, settings)
}

func (s *Server) handleUpdateCPASettings(w http.ResponseWriter, r *http.Request) {
	var settings model.CPASettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if settings.Endpoint == "" || settings.BearerToken == "" {
		jsonError(w, "endpoint and bearer_token are required", http.StatusBadRequest)
		return
	}

	if err := s.store.SaveCPASettings(r.Context(), &settings); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonOK(w, settings)
}
