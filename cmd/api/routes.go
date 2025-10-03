package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	{

		v1.GET(eventsPath, app.getAllEvents)
		v1.GET(eventIDPath, app.getEvent)

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.loginUser)
	}

	protectedRoutes := v1.Group("/")
	protectedRoutes.Use(app.middleware())

	{
		protectedRoutes.GET(eventAttendeesPath, app.getAttendees)
		protectedRoutes.POST(eventsPath, app.createEvent)
		protectedRoutes.PUT(eventIDPath, app.updateEvent)
		protectedRoutes.DELETE(eventIDPath, app.deleteEvent)

		protectedRoutes.POST(eventAttendeeIDPath, app.addAttendee)
		protectedRoutes.DELETE(eventAttendeeIDPath, app.deleteAttendee)

	}

	return g
}
