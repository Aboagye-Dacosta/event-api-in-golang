package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	{
		v1.POST(eventsPath, app.createEvent)
		v1.GET(eventsPath, app.getAllEvents)
		v1.GET(eventIDPath, app.getEvent)
		v1.PUT(eventIDPath, app.updateEvent)
		v1.DELETE(eventIDPath, app.deleteEvent)

		v1.POST("/auth/register", app.registerUser)
	}

	return g
}
