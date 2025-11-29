package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"share/devices"
)

func (s *Server) DevicePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "device.html")))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, map[string]string{
		"DeviceID": r.URL.Query().Get("id"),
	})
}

func (s *Server) ListDevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(s.registry.List())
}

func (s *Server) RegisterDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	device, err := s.registry.Upsert(payload.ID, payload.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(device)
}

func (s *Server) NotifyDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		DeviceID   string `json:"deviceId"`
		TransferID string `json:"transferId"`
		Token      string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	transfer, err := s.store.Authorize(payload.TransferID, payload.Token)
	if err != nil {
		writeTransferError(w, err)
		return
	}
	files := make([]devices.PendingFile, 0, len(transfer.Files))
	for _, f := range transfer.Files {
		files = append(files, devices.PendingFile{
			Name: f.Name,
			Mime: f.Mime,
			Size: f.Size,
		})
	}
	err = s.registry.Notify(payload.DeviceID, &devices.PendingTransfer{
		TransferID: transfer.ID,
		Token:      transfer.Token,
		Files:      files,
		SentAt:     time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DevicePendingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	deviceID := strings.TrimSpace(r.URL.Query().Get("id"))
	if deviceID == "" {
		http.Error(w, "missing device id", http.StatusBadRequest)
		return
	}
	pending, err := s.registry.Pending(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if pending == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(pending)
}

func (s *Server) ClearPendingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		DeviceID   string `json:"deviceId"`
		TransferID string `json:"transferId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if payload.DeviceID == "" {
		http.Error(w, "missing device id", http.StatusBadRequest)
		return
	}
	s.registry.Clear(payload.DeviceID, payload.TransferID)
	w.WriteHeader(http.StatusNoContent)
}
