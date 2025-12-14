package main

import (
	"log"
)

const (
	defaultDBPath  = "/usr/src/app/data/whoknows.db"
	defaultLogPath = "/usr/src/app/data/server.log"
)

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

	if _, err := openDatabase(defaultDBPath); err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer closeDatabase()

	go monitorUserCount(db)

	router := newRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
