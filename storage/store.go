package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	// ErrNotFound indicates that the provided transfer ID does not exist.
	ErrNotFound = errors.New("transfer not found")
	// ErrUnauthorized indicates that the provided token does not match the transfer.
	ErrUnauthorized = errors.New("invalid access token")
)

// Transfer holds information about an uploaded bundle.
type Transfer struct {
	ID        string
	Token     string
	Category  string
	PinHash   string
	Files     []StoredFile
	CreatedAt time.Time
}

// StoredFile represents a file inside a transfer bundle.
type StoredFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Mime string `json:"mime"`
	Size int64  `json:"size"`
	Path string
}

// FilePayload represents an uploaded file stream.
type FilePayload struct {
	Name    string
	MIME    string
	Content io.ReadCloser
}

// Store manages transfer metadata and the upload directory.
type Store struct {
	mu        sync.RWMutex
	dir       string
	transfers map[string]*Transfer
}

// NewStore creates a Store that persists files under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{
		dir:       dir,
		transfers: make(map[string]*Transfer),
	}, nil
}

// SaveFiles stores multiple files under a single transfer.
func (s *Store) SaveFiles(category, pin string, payloads []FilePayload) (*Transfer, error) {
	if len(payloads) == 0 {
		return nil, errors.New("no files provided")
	}
	defer func() {
		for _, payload := range payloads {
			if payload.Content != nil {
				payload.Content.Close()
			}
		}
	}()
	id := randomString(12)
	token := randomString(32)
	dir := filepath.Join(s.dir, id)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	transfer := &Transfer{
		ID:        id,
		Token:     token,
		Category:  category,
		PinHash:   HashPin(pin),
		CreatedAt: time.Now().UTC(),
	}

	for idx, payload := range payloads {
		if payload.Content == nil {
			continue
		}
		safe := sanitizeFilename(payload.Name)
		filename := fmt.Sprintf("%02d_%s", idx, safe)
		path := filepath.Join(dir, filename)
		out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			s.cleanupDir(dir)
			return nil, err
		}
		size, err := io.Copy(out, payload.Content)
		out.Close()
		if err != nil {
			_ = os.Remove(path)
			s.cleanupDir(dir)
			return nil, err
		}
		transfer.Files = append(transfer.Files, StoredFile{
			ID:   fmt.Sprintf("%s-%02d", id, idx),
			Name: payload.Name,
			Mime: payload.MIME,
			Size: size,
			Path: path,
		})
	}

	if len(transfer.Files) == 0 {
		s.cleanupDir(dir)
		return nil, errors.New("unable to store files")
	}

	s.mu.Lock()
	s.transfers[id] = transfer
	s.mu.Unlock()
	return transfer, nil
}

// Authorize returns the transfer if both the id and token are valid.
func (s *Store) Authorize(id, token string) (*Transfer, error) {
	s.mu.RLock()
	transfer, ok := s.transfers[id]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(transfer.Token)) != 1 {
		return nil, ErrUnauthorized
	}
	return transfer, nil
}

// Remove deletes the transfer metadata and file from disk.
func (s *Store) Remove(id string) {
	s.mu.Lock()
	_, ok := s.transfers[id]
	if ok {
		delete(s.transfers, id)
	}
	s.mu.Unlock()
	if ok {
		s.cleanupDir(filepath.Join(s.dir, id))
	}
}

// CleanupOlderThan removes any transfer that is older than ttl.
func (s *Store) CleanupOlderThan(ttl time.Duration) int {
	cutoff := time.Now().UTC().Add(-ttl)
	var removed int
	var paths []string

	s.mu.Lock()
	for id, transfer := range s.transfers {
		if transfer.CreatedAt.Before(cutoff) {
			delete(s.transfers, id)
			removed++
			paths = append(paths, id)
		}
	}
	s.mu.Unlock()
	for _, id := range paths {
		s.cleanupDir(filepath.Join(s.dir, id))
	}
	return removed
}

// StartCleanup periodically removes expired transfers until the context is done.
func (s *Store) StartCleanup(ctx context.Context, ttl, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.CleanupOlderThan(ttl)
		case <-ctx.Done():
			return
		}
	}
}

func sanitizeFilename(name string) string {
	base := filepath.Base(name)
	base = strings.Map(func(r rune) rune {
		switch {
		case r == '\\' || r == '/' || r == ':' || r == 0:
			return -1
		case r < 32:
			return -1
		}
		return r
	}, base)
	base = strings.TrimSpace(base)
	if base == "" {
		return "file.bin"
	}
	return base
}

func randomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)[:length]
}

// HashPin returns the SHA-256 hash of a trimmed PIN or empty string.
func HashPin(pin string) string {
	pin = strings.TrimSpace(pin)
	if pin == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(sum[:])
}

func (s *Store) cleanupDir(dir string) {
	_ = os.RemoveAll(dir)
}
