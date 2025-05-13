package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var (
	UploadDir        = "static/uploaded"
	UploadedFilePath string
	UploadedFileName string
	FileMime         string
)

func UploadPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "send.html")))
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll(UploadDir, os.ModePerm)
	destPath := filepath.Join(UploadDir, handler.Filename)
	f, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Save error", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	UploadedFilePath = destPath
	UploadedFileName = handler.Filename
	FileMime = handler.Header.Get("Content-Type")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
