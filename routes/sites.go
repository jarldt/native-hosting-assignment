package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"static-site-hosting/services"
	"strings"
)

/*
====================
ListSitesHandler

	Returns a list of deployed sites in JSON

====================
*/
func ListSitesHandler(db services.DatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sites, err := db.ListSites()
		if err != nil {
			fmt.Printf("Error listing sites: %v\n", err)
			http.Error(w, "Failed to list sites", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sites)
	}
}

/*
====================
DeleteSiteHandler

	Deletes a deployed site from the DB and storage

====================
*/
func DeleteSiteHandler(db services.DatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract siteName from URL path
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 2 || parts[0] != "delete" {
			http.Error(w, "Invalid delete path", http.StatusBadRequest)
			return
		}
		siteName := parts[1]

		// Validate siteName
		if siteName == "" || strings.Contains(siteName, "..") || strings.Contains(siteName, "/") {
			http.Error(w, "Invalid site name", http.StatusBadRequest)
			return
		}

		// Delete site folder
		sitePath := path.Join("deployments", siteName)
		if _, err := os.Stat(sitePath); os.IsNotExist(err) {
			http.Error(w, "Site not found", http.StatusNotFound)
			return
		}

		if err := os.RemoveAll(sitePath); err != nil {
			fmt.Printf("Error deleting site files: %v\n", err)
			http.Error(w, "Failed to delete site files", http.StatusInternalServerError)
			return
		}

		// Delete from DB
		err := db.DeleteSiteByName(siteName)
		if err != nil {
			fmt.Printf("Error deleting site from DB: %v\n", err)
			http.Error(w, "Failed to delete site from database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
