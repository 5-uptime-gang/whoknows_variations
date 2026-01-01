package main

import (
	"log"
)

const (
	defaultLogPath = "/usr/src/app/data/server.log"
)

// @title WhoKnows Variations API
// @version 1.0
// @description Routes for pages, auth, search, metrics, and weather.
// @BasePath /
// @host localhost:8080
// @schemes http
func main() {
	logFile, err := setupLogging(defaultLogPath)
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}()

	database, err := openDatabase()
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer closeDatabase()

	if err := InitDB(database); err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	go monitorUserCount(db)

	router := newRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
