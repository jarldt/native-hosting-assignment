package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"static-site-hosting/services"
)

// DeploymentLogsHandler returns recent deployment logs
func DeploymentLogsHandler(db services.DatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// TODO: Validate site name from path
		site := r.PathValue("site")
		logs, err := db.GetDeploymentLogs(site)
		if err != nil {
			fmt.Printf("Error getting deployment logs for site '%s': %v\n", site, err)
			http.Error(w, "Failed to get deployment logs", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}
