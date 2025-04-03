package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	handler "tailwind-converter/api"
)

func main() {
	// Define command line parameters
	port := 8080

	// Create a new router
	mux := http.NewServeMux()

	// Create static directory if it doesn't exist
	ensureStaticDirectory()

	// Register API handlers
	mux.HandleFunc("/convert", handler.Handler)

	// Serve static files (for assets referenced by index.html)
	staticHandler := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set proper content types based on file extensions
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}
		staticHandler.ServeHTTP(w, r)
	})))

	// Also serve script.js and other root assets directly
	mux.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "static/script.js")
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

// ensureStaticDirectory creates the static directory and copies files from public if needed
func ensureStaticDirectory() {
	// Check if static directory exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		log.Println("Static directory not found, creating from public directory...")

		// Create static directory
		err := os.MkdirAll("static", 0o755)
		if err != nil {
			log.Fatalf("Failed to create static directory: %v", err)
		}

		// Copy files from public to static
		err = copyDir("public", "static")
		if err != nil {
			log.Fatalf("Failed to copy files from public to static: %v", err)
		}

		log.Println("Static directory created successfully")
	}
}

// copyDir copies files from source directory to destination directory
func copyDir(src, dst string) error {
	// Get file info
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if source is directory
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// Create destination directory
	err = os.MkdirAll(dst, info.Mode())
	if err != nil {
		return err
	}

	// Read directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursive copy for directories
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write to destination file
	return os.WriteFile(dst, data, 0o644)
}
