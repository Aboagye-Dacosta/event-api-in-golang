package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	{
		g.GET("/swagger/*any", func(c *gin.Context) {

			log.Println("swagger", c.Request.RequestURI)

			if c.Request.RequestURI == "/swagger" || c.Request.RequestURI == "/swagger/" {
				c.Redirect(302, "/swagger/index.html")
			}

			ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
			// c.JSON(http.StatusOK, gin.H{"message": "swagger"})
		})

	}

	return g
}
