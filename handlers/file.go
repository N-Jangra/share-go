package handlers

import (
	"net/http"
)

func ServeFileHandler(w http.ResponseWriter, r *http.Request) {
	if UploadedFileName == "" {
		http.Error(w, "No file uploaded", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+UploadedFileName+"\"")
	w.Header().Set("Content-Type", FileMime)
	http.ServeFile(w, r, UploadedFilePath)
}
