package api

import "net/http"

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.ListAccountsWithQuota(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	syncLog, _ := s.store.GetLatestSyncLog(r.Context())

	jsonOK(w, map[string]any{
		"accounts": accounts,
		"cpa_sync": syncLog,
	})
}
