package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"share/devices"
	"share/storage"
)

// Server groups all HTTP handlers together with their dependencies.
type Server struct {
	store     *storage.Store
	registry  *devices.Registry
	pinSecret []byte
}

// NewServer builds a handler server with the provided storage backend.
func NewServer(store *storage.Store, registry *devices.Registry) *Server {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		panic(err)
	}
	return &Server{
		store:     store,
		registry:  registry,
		pinSecret: secret,
	}
}

func (s *Server) pinCookieName(id string) string {
	return "pin-" + id
}

func (s *Server) pinCookieValue(t *storage.Transfer) string {
	mac := hmac.New(sha256.New, s.pinSecret)
	mac.Write([]byte(t.PinHash))
	mac.Write([]byte(t.ID))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Server) hasPinAccess(r *http.Request, t *storage.Transfer) bool {
	if t.PinHash == "" {
		return true
	}
	cookie, err := r.Cookie(s.pinCookieName(t.ID))
	if err != nil {
		return false
	}
	return cookie.Value == s.pinCookieValue(t)
}

func (s *Server) grantPinAccess(w http.ResponseWriter, t *storage.Transfer) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.pinCookieName(t.ID),
		Value:    s.pinCookieValue(t),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int((2 * time.Hour).Seconds()),
	})
}

func (s *Server) validatePin(input string, t *storage.Transfer) bool {
	if t.PinHash == "" {
		return true
	}
	return storage.HashPin(input) == t.PinHash
}
