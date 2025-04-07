package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Index function for serving static files or redirecting to index.html
func Index(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log the request path
	log.Printf("Handling request for path: %s", r.URL.Path)

	// If the request is for root, serve index.html
	if r.URL.Path == "/" {
		serveIndexHTML(w, r)
		return
	}

	// Check if this is a request for a static file
	if isStaticFileRequest(r.URL.Path) {
		serveStaticFile(w, r)
		return
	}

	// Default to serving index.html for all other paths
	serveIndexHTML(w, r)
}

// isStaticFileRequest checks if the request is for a static file
func isStaticFileRequest(path string) bool {
	extensions := []string{".js", ".css", ".png", ".jpg", ".gif", ".svg", ".ico", ".html"}
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

// serveStaticFile serves a static file from the public directory
func serveStaticFile(w http.ResponseWriter, r *http.Request) {
	// Get the file path from the URL
	filePath := r.URL.Path

	// If the path starts with /static/, remove it
	if strings.HasPrefix(filePath, "/static/") {
		filePath = strings.TrimPrefix(filePath, "/static/")
	}

	// Construct the full path to the file in the public directory
	fullPath := filepath.Join("public", filePath)

	// Log the file path
	log.Printf("Attempting to serve static file: %s", fullPath)

	// Check if the file exists
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		// Try alternative paths in Vercel environment
		if os.Getenv("VERCEL") != "" {
			alternatives := []string{
				filepath.Join("/var/task/public", filePath),
				filepath.Join("/public", filePath),
			}

			for _, alt := range alternatives {
				if _, err := os.Stat(alt); err == nil {
					fullPath = alt
					log.Printf("Found file at alternative path: %s", fullPath)
					break
				}
			}
		}
	}

	// Check if the file exists at the resolved path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("File not found: %s", fullPath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set the appropriate content type based on file extension
	setContentType(w, fullPath)

	// Serve the file
	http.ServeFile(w, r, fullPath)
}

// serveIndexHTML serves the index.html file
func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	// Define possible index.html paths
	possiblePaths := []string{
		"public/index.html",
		"/var/task/public/index.html",
		"/public/index.html",
	}

	// Try to find and serve the index.html file
	for _, path := range possiblePaths {
		content, err := os.ReadFile(path)
		if err == nil {
			log.Printf("Serving index.html from: %s", path)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(content)
			return
		}
	}

	// If we couldn't find the index.html file, return an error
	log.Printf("Error: Could not find index.html in any location")
	http.Error(w, "Unable to find index.html file", http.StatusInternalServerError)
}

// setContentType sets the appropriate Content-Type header based on file extension
func setContentType(w http.ResponseWriter, path string) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}
