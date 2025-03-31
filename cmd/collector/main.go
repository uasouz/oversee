package main

import (
	"log"
	"oversee/collector/logsapi"
	"oversee/collector/persistence/sqlite"
)

func main() {
	sqlitePersistence, err := sqlite.NewSQLitePersistence("test.db")

	if err != nil {
		log.Fatal("Failed to initialize persistence")
	}

	collectorApi := logsapi.NewLogsAPI(sqlitePersistence)

	collectorApi.Serve()
}
