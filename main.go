package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"share/handlers"
	"share/utils"
)

func main() {
	// Ensure upload directory exists
	err := createUploadDir()
	if err != nil {
		log.Fatalf("Error creating upload directory: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/upload", handlers.UploadPage)
	http.HandleFunc("/uploadFile", handlers.UploadFileHandler)
	http.HandleFunc("/meta", handlers.FileMetaHandler)
	http.HandleFunc("/file", handlers.ServeFileHandler)
	http.HandleFunc("/incoming", handlers.IncomingHandler)
	http.HandleFunc("/accept", handlers.AcceptHandler)
	http.HandleFunc("/decline", handlers.DeclineHandler)

	port := "8080"
	ips := utils.GetAllLocalIPs()
	for _, ip := range ips {
		fmt.Printf("Server running at: http://%s:%s\n", ip, port)
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createUploadDir() error {
	return ensureDir(filepath.Join("static", "uploaded"))
}

func ensureDir(dir string) error {
	return utils.EnsureDir(dir)
}
