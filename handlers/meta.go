package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func FileMetaHandler(w http.ResponseWriter, r *http.Request) {
	if UploadedFileName == "" {
		http.Error(w, "No file uploaded", http.StatusNotFound)
		return
	}

	fi, _ := os.Stat(UploadedFilePath)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"name": "%s", "size": %d, "type": "%s"}`, UploadedFileName, fi.Size(), FileMime)
}
