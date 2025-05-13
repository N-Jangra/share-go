package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"share/utils"
)

type IncomingData struct {
	Name   string
	SizeMB float64
	Type   string
}

func IncomingHandler(w http.ResponseWriter, r *http.Request) {
	meta, err := utils.GetFileMeta(getSenderURL())
	if err != nil {
		fmt.Fprint(w, "<p>No file available from sender yet.</p>")
		return
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "receive.html")))
	data := IncomingData{
		Name:   meta.Name,
		SizeMB: float64(meta.Size) / (1024 * 1024),
		Type:   meta.Type,
	}
	tmpl.Execute(w, data)
}

func AcceptHandler(w http.ResponseWriter, r *http.Request) {
	err := downloadFile(w)
	if err != nil {
		fmt.Fprintf(w, "<p>Download failed: %s</p>", err)
	}
}

func DeclineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<p>Transfer declined.</p>")
}

func downloadFile(w http.ResponseWriter) error {
	resp, err := http.Get(getSenderURL() + "/file")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	meta, err := utils.GetFileMeta(getSenderURL())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+meta.Name+"\"")
	w.Header().Set("Content-Type", meta.Type)
	_, err = io.Copy(w, resp.Body)
	return err
}

func getSenderURL() string {
	ip := utils.GetLocalIP()
	return fmt.Sprintf("http://%s:8080", ip)
}
