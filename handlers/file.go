package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (s *Server) ServeFileHandler(w http.ResponseWriter, r *http.Request) {
	transfer, err := s.transferFromRequest(r)
	if err != nil {
		writeTransferError(w, err)
		return
	}
	if transfer.PinHash != "" && !s.hasPinAccess(r, transfer) {
		http.Error(w, "pin required", http.StatusForbidden)
		return
	}
	fileID := r.URL.Query().Get("file")
	stored, err := findFile(transfer, fileID)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(stored.Path)
	if err != nil {
		http.Error(w, "file unavailable", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stored.Name))
	w.Header().Set("Content-Type", stored.Mime)
	w.Header().Set("Content-Length", strconv.FormatInt(stored.Size, 10))
	if _, err := io.Copy(w, f); err != nil {
		return
	}
}
