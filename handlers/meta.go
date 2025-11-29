package handlers

import (
	"encoding/json"
	"net/http"
)

func (s *Server) FileMetaHandler(w http.ResponseWriter, r *http.Request) {
	transfer, err := s.transferFromRequest(r)
	if err != nil {
		writeTransferError(w, err)
		return
	}
	if transfer.PinHash != "" && !s.hasPinAccess(r, transfer) {
		http.Error(w, "pin required", http.StatusForbidden)
		return
	}

	files := make([]map[string]interface{}, 0, len(transfer.Files))
	for _, f := range transfer.Files {
		files = append(files, map[string]interface{}{
			"id":   f.ID,
			"name": f.Name,
			"mime": f.Mime,
			"size": f.Size,
		})
	}
	resp := map[string]interface{}{
		"category":    transfer.Category,
		"requiresPin": transfer.PinHash != "",
		"files":       files,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(resp)
}
