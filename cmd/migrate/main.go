package main

import (
	// The lines `"database/sql"`, `"log"`, and `"os"` are import statements in Go.
	"database/sql"
	"log"
	"os"

	// These lines are import statements in Go that are used to import external packages into the current Go program. Here's what each import statement is doing:
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// The line `fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")` is opening a file source for migrations.
	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)

	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid direction use: up or down")
	}
}
