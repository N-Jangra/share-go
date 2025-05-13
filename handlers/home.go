package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "main.html")))
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}
