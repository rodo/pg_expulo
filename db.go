package main

import (
	"os"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

func get_dsn(host string, port string, user string, pass string, db string, version string) string {

	appname := "expulo_" + version

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name=%s sslmode=disable", host, port, user, pass, db, appname)
	return dsn
}

func connectDb(connectionString string) *sql.DB {
	// Connect to the database source
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Error connecting to the database source:", err)
		os.Exit(1)
	}

	// Ensure the connection is up and running
	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
