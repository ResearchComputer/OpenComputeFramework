package access

import (
	"context"
	"log"
	"ocfcore/internal/database"
)

type DBClient struct {
	database.Client
}

var dbClient *DBClient

func NewDBClient() *DBClient {
	if dbClient == nil {
		dbClient, err := database.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
		}
		// Run the auto migration tool.
		if err := dbClient.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}
	return dbClient
}
