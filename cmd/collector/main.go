package main

import (
	"log"
	"oversee/collector/audit"
	"oversee/collector/graphql"
	"oversee/collector/persistence/sqlite"
)

func main() {
	sqlitePersistence, err := sqlite.NewSQLitePersistence("test.db")

	if err != nil {
		log.Fatal("Failed to initialize persistence")
	}

	collectorApi := audit.NewLogsIngestionAPI(sqlitePersistence)

	go func() {
		searchService := audit.NewSearchService(sqlitePersistence)
		gqlServer := graphql.NewGraphqlAPIServer(searchService)
		gqlServer.Start()
	}()

	collectorApi.Serve()
}
