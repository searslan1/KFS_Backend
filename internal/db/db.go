package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}

	fmt.Println("Successfully connected to the database")
}

func RunMigrations(filePath string) {
	migration, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading migration file: %q", err)
	}

	_, err = DB.Exec(string(migration))
	if err != nil {
		log.Fatalf("Error executing migration: %q", err)
	}

	fmt.Println("Migration executed successfully")
}
