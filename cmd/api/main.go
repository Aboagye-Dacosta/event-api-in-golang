package main

import (
	"database/sql"
	"first-rest-api/internal/database"
	"first-rest-api/internal/env"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port       int
	jwtSecrete string
	models     database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "../../data.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModel(db)

	app := &application{
		port:       env.GetEnvInt("P0RT", 8080),
		jwtSecrete: env.GetEnvString("JWT_SECRETE", "some-secrete-123456"),
		models:     models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

}
