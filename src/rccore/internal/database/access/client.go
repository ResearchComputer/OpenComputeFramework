package access

import (
	"context"
	"log"

	"entgo.io/ent/examples/fs/ent"
)

type DBClient struct {
	*ent.Client
}

var dbClient *DBClient

func NewDBClient() *DBClient {
	if dbClient == nil {
		dbClient, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
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
