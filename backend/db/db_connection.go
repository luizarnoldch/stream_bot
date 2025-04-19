package db

import (
	"database/sql"
	"log"
)

// GetConn sets up and returns a database connection.
func GetConn(dsn string, maxConns int) *sql.DB {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DbConn() Failed to connect to database, dsn=%s, err=%s", dsn, err.Error())
		return dbConn
	}

	dbConn.SetMaxOpenConns(maxConns)     // Maximum number of open connections to the database
	dbConn.SetMaxIdleConns(maxConns / 2) // Maximum number of idle connections to the database

	return dbConn
}
