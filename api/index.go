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

	// Log the current directory for debugging
	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)

	// Check if running in Vercel environment
	inVercel := os.Getenv("VERCEL") != ""
	log.Printf("Running in Vercel environment: %v", inVercel)

	// Define possible public directory paths
	possiblePaths := []string{
		"/var/task/public",
		filepath.Join(cwd, "public"),
		filepath.Join(cwd, "../public"),
		"/public",
		"./public",
		"../public",
		"/home/site/public",
		"/vercel/path0/public",
	}

	// Add more debug info
	log.Printf("Checking the following paths for index.html: %v", possiblePaths)

	// Find index.html in one of the possible public directories
	var indexPath string
	var found bool

	for _, dir := range possiblePaths {
		path := filepath.Join(dir, "index.html")
		_, err := os.Stat(path)
		if err == nil {
			indexPath = path
			found = true
			log.Printf("Found index.html at: %s", path)
			break
		} else {
			log.Printf("Could not find index.html at %s: %v", path, err)
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

		// Try to list the contents of /var/task directly
		if taskFiles, err := os.ReadDir("/var/task"); err == nil {
			log.Printf("Files in /var/task directory:")
			for _, file := range taskFiles {
				log.Printf("- %s (is dir: %v)", file.Name(), file.IsDir())
			}
		}

		// Check if the public directory exists in /var/task
		if publicInfo, err := os.Stat("/var/task/public"); err == nil {
			log.Printf("/var/task/public exists, is directory: %v", publicInfo.IsDir())
			// List contents of /var/task/public
			if publicFiles, err := os.ReadDir("/var/task/public"); err == nil {
				log.Printf("Files in /var/task/public directory:")
				for _, file := range publicFiles {
					log.Printf("- %s (is dir: %v)", file.Name(), file.IsDir())
				}
			}
		} else {
			log.Printf("Error accessing /var/task/public: %v", err)
		}

		http.Error(w, "Unable to find static files. Please ensure public directory exists.", http.StatusInternalServerError)
		return
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
