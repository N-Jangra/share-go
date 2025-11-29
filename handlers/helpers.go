package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"share/storage"
)

func (s *Server) transferFromRequest(r *http.Request) (*storage.Transfer, error) {
	id := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")
	if id == "" || token == "" {
		return nil, errors.New("missing id or token")
	}
	return s.store.Authorize(id, token)
}

func writeTransferError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		http.Error(w, "transfer not found", http.StatusNotFound)
	case errors.Is(err, storage.ErrUnauthorized):
		http.Error(w, "invalid token", http.StatusForbidden)
	default:
		http.Error(w, "unable to process request", http.StatusBadRequest)
	}
}

func findFile(transfer *storage.Transfer, fileID string) (*storage.StoredFile, error) {
	if fileID == "" {
		return nil, fmt.Errorf("missing file id")
	}
	for i := range transfer.Files {
		if transfer.Files[i].ID == fileID {
			return &transfer.Files[i], nil
		}
	}
	return nil, fmt.Errorf("file not found")
}
