package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// IndexHandler function for serving the index.html file
func IndexHandler(w http.ResponseWriter, r *http.Request) {
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

	// Try multiple possible paths for static directory
	staticDirs := []string{
		"./static",             // Local development
		"/var/task/static",     // Vercel
		"/home/site/static",    // Another possible Vercel path
		"/vercel/path0/static", // Another possible Vercel path
	}

	var indexPath string
	var staticDir string

	// Find the first valid static directory
	for _, dir := range staticDirs {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			testPath := filepath.Join(dir, "index.html")
			if _, err := os.Stat(testPath); !os.IsNotExist(err) {
				indexPath = testPath
				staticDir = dir
				log.Printf("Found static directory at: %s", staticDir)
				break
			}
		}
	}

	// If no valid path found
	if indexPath == "" {
		log.Printf("Error: Could not find static files in any of the expected locations")

		// List the current directory to help debug
		files, _ := filepath.Glob("*")
		log.Printf("Files in current directory: %v", files)

		http.Error(w, "Unable to find static files", http.StatusInternalServerError)
		return
	}

	// Read the index.html file
	content, err := ioutil.ReadFile(indexPath)
	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		http.Error(w, "Unable to read index file", http.StatusInternalServerError)
		return
	}

	// Set content type and serve the file
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
