package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed init.sql
var embedded embed.FS

func OpenDatabase(path string) *sql.DB {
	dsn := fmt.Sprintf("file:%s", path)
	db, err := connect(dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = setup(db)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return db, err
	}
	return db, nil
}

func ping(db sql.DB) bool {
	rows, err := db.Query("SELECT 1")
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	var i int
	rows.Scan(&i)
	log.Println(i)
	return i == 1
}

func setup(db *sql.DB) error {
	_, err := db.Query("SELECT * FROM listing;")
	if err == nil {
		log.Println("Database already initialized")
		return nil
	}

	bytes, err := embedded.ReadFile("init.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(bytes))
	if err != nil {
		return err
	}
	return nil
}
