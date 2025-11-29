package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
)

type incomingFile struct {
	ID     string
	Name   string
	Mime   string
	SizeMB float64
}

type IncomingPageData struct {
	ID          string
	Token       string
	Category    string
	RequiresPin bool
	NeedsPin    bool
	PinError    string
	Files       []incomingFile
}

func (s *Server) IncomingHandler(w http.ResponseWriter, r *http.Request) {
	transfer, err := s.transferFromRequest(r)
	if err != nil {
		writeTransferError(w, err)
		return
	}

	data := IncomingPageData{
		ID:          transfer.ID,
		Token:       transfer.Token,
		Category:    categoryLabel(transfer.Category),
		RequiresPin: transfer.PinHash != "",
	}
	hasAccess := s.hasPinAccess(r, transfer)
	if r.Method == http.MethodPost && transfer.PinHash != "" && !hasAccess {
		if err := r.ParseForm(); err == nil {
			pin := r.FormValue("pin")
			if s.validatePin(pin, transfer) {
				s.grantPinAccess(w, transfer)
				http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
				return
			}
			data.PinError = "Incorrect PIN"
		}
	}
	if transfer.PinHash != "" && !s.hasPinAccess(r, transfer) {
		data.NeedsPin = true
	} else {
		for _, f := range transfer.Files {
			data.Files = append(data.Files, incomingFile{
				ID:     f.ID,
				Name:   f.Name,
				Mime:   f.Mime,
				SizeMB: float64(f.Size) / (1024 * 1024),
			})
		}
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "receive.html")))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}

func (s *Server) AcceptHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")
	fileID := r.URL.Query().Get("file")
	if id == "" || token == "" {
		http.Error(w, "missing transfer information", http.StatusBadRequest)
		return
	}
	params := url.Values{}
	params.Set("id", id)
	params.Set("token", token)
	if fileID != "" {
		params.Set("file", fileID)
	}
	target := fmt.Sprintf("/file?%s", params.Encode())
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (s *Server) DeclineHandler(w http.ResponseWriter, r *http.Request) {
	transfer, err := s.transferFromRequest(r)
	if err != nil {
		writeTransferError(w, err)
		return
	}
	if transfer.PinHash != "" && !s.hasPinAccess(r, transfer) {
		http.Error(w, "pin required", http.StatusForbidden)
		return
	}
	s.store.Remove(transfer.ID)
	s.registry.ClearByTransfer(transfer.ID)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<div class='card'><p>Transfer declined and removed.</p><a href='/' class='button btn-secondary'>Return Home</a></div>")
}
