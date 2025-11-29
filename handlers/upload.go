package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"text/template"

	"share/storage"
)

const maxUploadSize = 25 << 20 // 25 MB

func (s *Server) UploadPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "send.html")))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, nil)
}

func (s *Server) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "failed to parse upload", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "please attach at least one file", http.StatusBadRequest)
		return
	}
	var payloads []storage.FilePayload
	for _, header := range files {
		f, err := header.Open()
		if err != nil {
			http.Error(w, "unable to read file", http.StatusBadRequest)
			return
		}
		payloads = append(payloads, storage.FilePayload{
			Name:    header.Filename,
			MIME:    header.Header.Get("Content-Type"),
			Content: f,
		})
	}
	category := normalizeCategory(r.FormValue("category"))
	var pin string
	if strings.TrimSpace(r.FormValue("pin")) != "" {
		pin = r.FormValue("pin")
	}
	transfer, err := s.store.SaveFiles(category, pin, payloads)
	if err != nil {
		http.Error(w, "unable to store files", http.StatusInternalServerError)
		return
	}

	scheme := requestScheme(r)
	shareURL := fmt.Sprintf("%s://%s/incoming?id=%s&token=%s", scheme, r.Host, transfer.ID, transfer.Token)
	filesData := make([]shareFile, 0, len(transfer.Files))
	for _, f := range transfer.Files {
		filesData = append(filesData, shareFile{
			Name:      f.Name,
			Mime:      f.Mime,
			SizeMB:    float64(f.Size) / (1024 * 1024),
			DirectURL: fmt.Sprintf("%s://%s/file?id=%s&token=%s&file=%s", scheme, r.Host, transfer.ID, transfer.Token, url.QueryEscape(f.ID)),
		})
	}
	data := sharePageData{
		ShareLink:   shareURL,
		Category:    categoryLabel(category),
		RequiresPin: transfer.PinHash != "",
		TransferID:  transfer.ID,
		Token:       transfer.Token,
		Files:       filesData,
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "share.html")))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}

type shareFile struct {
	Name      string
	Mime      string
	SizeMB    float64
	DirectURL string
}

type sharePageData struct {
	ShareLink   string
	Category    string
	RequiresPin bool
	TransferID  string
	Token       string
	Files       []shareFile
}

func normalizeCategory(input string) string {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "audio", "audios":
		return "audio"
	case "photo", "photos", "images":
		return "photos"
	case "video", "videos":
		return "videos"
	case "document", "documents", "docs":
		return "documents"
	case "contact", "contacts":
		return "contacts"
	default:
		return "any"
	}
}

func categoryLabel(cat string) string {
	switch cat {
	case "audio":
		return "Audio"
	case "photos":
		return "Photos"
	case "videos":
		return "Videos"
	case "documents":
		return "Documents"
	case "contacts":
		return "Contacts"
	default:
		return "Any File"
	}
}

func requestScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
