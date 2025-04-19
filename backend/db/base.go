package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	instance *Queries
	once     sync.Once
)

// Get returns the database.  Optionally accepts a `uri` string to connect to.
// If no DSN string is passed, looks for and connects to the `uri` environmental variable.
//
// Example uri:
// "postgresql://postgres:postgres@localhost:5434/postgres?pool_max_conns=10"
func GetPSQLClient(uri string) *Queries {
	once.Do(func() {
		pool, err := pgxpool.New(context.Background(), uri)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		instance = New(pool)
	})
	return instance
}
