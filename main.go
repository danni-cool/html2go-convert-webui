package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	handler "html2go-converter/api"
)

func main() {
	// Define command line parameters
	portPtr := flag.Int("port", 8080, "端口号")
	flag.Parse()
	port := *portPtr

	// Create a new router
	mux := http.NewServeMux()

	// Register API handlers
	mux.HandleFunc("/convert", handler.Handler)

	// Serve static files from public directory
	publicHandler := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set proper content types based on file extensions
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}
		publicHandler.ServeHTTP(w, r)
	})))

	// Also serve script.js and other root assets directly
	mux.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "public/script.js")
	})

	// Handle root path last
	mux.HandleFunc("/", handler.Index)

	// Configure the HTTP server
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port
	log.Printf("Starting server on http://localhost:%d", actualPort)
	log.Fatal(server.Serve(listener))
}
