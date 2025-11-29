package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"share/devices"
	"share/handlers"
	"share/storage"
	"share/utils"
)

const (
	port            = "8080"
	uploadsDir      = "uploads"
	transferTTL     = 2 * time.Hour
	cleanupInterval = 15 * time.Minute
)

func main() {
	store, err := storage.NewStore(filepath.Clean(uploadsDir))
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	registry := devices.NewRegistry(50)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go store.StartCleanup(ctx, transferTTL, cleanupInterval)

	server := handlers.NewServer(store, registry)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", server.HomeHandler)
	http.HandleFunc("/upload", server.UploadPage)
	http.HandleFunc("/uploadFile", server.UploadFileHandler)
	http.HandleFunc("/meta", server.FileMetaHandler)
	http.HandleFunc("/file", server.ServeFileHandler)
	http.HandleFunc("/incoming", server.IncomingHandler)
	http.HandleFunc("/accept", server.AcceptHandler)
	http.HandleFunc("/decline", server.DeclineHandler)
	http.HandleFunc("/device", server.DevicePage)
	http.HandleFunc("/api/devices", server.ListDevicesHandler)
	http.HandleFunc("/api/devices/register", server.RegisterDeviceHandler)
	http.HandleFunc("/api/devices/notify", server.NotifyDeviceHandler)
	http.HandleFunc("/api/devices/pending", server.DevicePendingHandler)
	http.HandleFunc("/api/devices/clear", server.ClearPendingHandler)

	ips := utils.GetAllLocalIPs()
	if len(ips) == 0 {
		fmt.Printf("Server running at: http://localhost:%s\n", port)
	}
	for _, ip := range ips {
		fmt.Printf("Server running at: http://%s:%s\n", ip, port)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
