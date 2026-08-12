package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]time.Time),
	}
}

func (ss *SessionStore) Create() string {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	ss.mu.Lock()
	ss.sessions[token] = time.Now().Add(24 * time.Hour)
	ss.mu.Unlock()

	return token
}

func (ss *SessionStore) Valid(token string) bool {
	ss.mu.RLock()
	expires, ok := ss.sessions[token]
	ss.mu.RUnlock()

	if !ok {
		return false
	}
	if time.Now().After(expires) {
		ss.mu.Lock()
		delete(ss.sessions, token)
		ss.mu.Unlock()
		return false
	}
	return true
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || !s.sessions.Valid(cookie.Value) {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
