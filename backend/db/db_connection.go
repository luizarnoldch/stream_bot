package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // registra el driver de PostgreSQL
)

// GetConn sets up and returns a database connection.
func GetConn(dsn string, maxConns int) *sql.DB {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DbConn() Failed to connect to database, dsn=%s, err=%s", dsn, err.Error())
		return dbConn
	}

	dbConn.SetMaxOpenConns(maxConns)
	dbConn.SetMaxIdleConns(maxConns / 2)

	return dbConn
}
