package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"static-site-hosting/middleware"
	"static-site-hosting/routes"
	"static-site-hosting/services"
)

func main() {
	// Setup and connect to the database
	err := os.Remove("./db/database.db")
	if err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("failed to remove database file: %v", err))
	}

	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	dbService := services.NewDatabaseService(db)
	err = dbService.Initialize()
	if err != nil {
		panic(err)
	}

	// Create a sample table for demonstration purposes
	/* 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS example (id INTEGER PRIMARY KEY, name TEXT)")
	   	if err != nil {
	   		log.Fatalf("Error creating table: %v", err)
	   	} */

	mux := http.NewServeMux()
	// mux.HandleFunc("/hello-world", routes.HelloWorldHandler)
	mux.HandleFunc("/deploy", routes.DeployHandler(dbService))
	mux.HandleFunc("/sites", routes.ListSitesHandler(dbService))
	mux.HandleFunc("/delete/", routes.DeleteSiteHandler(dbService))
	mux.HandleFunc("/sites/{site}/{path...}", routes.ServeSiteHandler)
	mux.HandleFunc("/logs/{site}", routes.DeploymentLogsHandler(dbService))

	wrappedMux := middleware.LoggingMiddleware(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}
