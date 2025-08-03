package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/*
====================
ServeSiteHandler

	Serves static files from deployments/{siteName}/

====================
*/
func ServeSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteName := r.PathValue("site")
	if siteName == "" {
		http.NotFound(w, r)
		return
	}

	// Get the remaining path after the site name
	remainingPath := r.PathValue("path")
	if remainingPath == "" {
		remainingPath = "index.html"
	}

	// Prevent directory traversal
	if strings.Contains(remainingPath, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("deployments", siteName, remainingPath)

	info, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// If it's a directory, try to serve index.html from it
	if info.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
		info, err = os.Stat(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	http.ServeContent(w, r, filePath, info.ModTime(), file)
}
