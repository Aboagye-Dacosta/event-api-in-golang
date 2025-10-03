package main

import (
	"database/sql"
	_ "first-rest-api/docs"
	"first-rest-api/internal/database"
	"first-rest-api/internal/env"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// @title           Event Rest API
// @version         1.0
// @description     This is a simple REST API built with Go. It allows users to create, read, update, and delete events. The API is built using the Gin framework.
//
// @contact.name    Solomon Aboagye
// @contact.url     https://github.com/aboagye-dacosta
// @contact.email   dacostaaboagyesolomon@gmail.com
//
// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @description  Type "Bearer Token" in the format **Bearer {token}** to authenticate the request.
//
// @host      localhost:8080
// @BasePath  /api/v1


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
