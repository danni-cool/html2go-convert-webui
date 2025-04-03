package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Index function for serving the index.html file
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

	// Define possible public directory paths
	possiblePaths := []string{
		"./public",
		"../public",
		"/var/task/public",
		"/home/site/public",
		"/vercel/path0/public",
	}

	// Find index.html in one of the possible public directories
	var indexPath string
	var found bool

	for _, dir := range possiblePaths {
		path := filepath.Join(dir, "index.html")
		if _, err := os.Stat(path); err == nil {
			indexPath = path
			found = true
			log.Printf("Found index.html at: %s", path)
			break
		}
	}

	// If not found in predefined paths, try to find it in the current directory structure
	if !found {
		log.Printf("Error: Could not find index.html in expected locations")

		// Log directory contents for debugging
		files, _ := filepath.Glob("*")
		log.Printf("Files in current directory: %v", files)

		parentFiles, _ := filepath.Glob("../*")
		log.Printf("Files in parent directory: %v", parentFiles)

		// In Vercel environment, we know the structure better
		if _, exists := os.LookupEnv("VERCEL"); exists {
			// Try the hardcoded Vercel public path
			indexPath = "/var/task/public/index.html"
			if _, err := os.Stat(indexPath); err == nil {
				found = true
				log.Printf("Found index.html in Vercel path: %s", indexPath)
			}
		}

		if !found {
			http.Error(w, "Unable to find static files. Please ensure public directory exists.", http.StatusInternalServerError)
			return
		}
	}

	// Read the index.html file
	content, err := os.ReadFile(indexPath)
	if err != nil {
		log.Printf("Error reading index.html at %s: %v", indexPath, err)
		http.Error(w, "Unable to read index file", http.StatusInternalServerError)
		return
	}

	// Set content type and serve the file
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
