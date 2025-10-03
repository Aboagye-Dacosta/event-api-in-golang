package main

import (
	"first-rest-api/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) addAttendee(c *gin.Context) {
	// Get the event ID from the URL parameters
	eventID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		app.badRequest(c, err)
		return
	}
	// Get the user ID from the URL parameters
	userID, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		app.badRequest(c, err)
		return
	}

	// Check if the user is already attending the event
	_, err = app.models.Attendees.GetAttendeeByUserIDAndEventID(userID, eventID)

	if err == nil {
		app.editConflict(c)
		return
	}

	err = app.models.Attendees.Insert(&database.Attendees{
		UserId:  userID,
		EventId: eventID,
	})

	if err != nil {
		app.serverError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendee added successfully"})

	// Add the attendee to the event

}

func (m *application) getAttendees(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		m.badRequest(c, err)
		return
	}

	attendees, err := m.models.Attendees.GetAllAttendeesByEventID(eventID)

	if err != nil {
		m.serverError(c, err)
		return
	}

	users := make([]database.User, len(attendees))

	for i, attendee := range attendees {
		user, err := m.models.Users.Get(attendee.UserId)

		if err != nil {
			m.serverError(c, err)
			return
		}

		users[i] = *user
	}

	c.JSON(http.StatusOK, users)
}

func (app *application) deleteAttendee(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		app.badRequest(c, err)
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		app.badRequest(c, err)
		return
	}

	attendee, err := app.models.Attendees.GetAttendeeByUserIDAndEventID(userID, eventID)

	if err != nil {
		app.notFound(c)
		return
	}

	err = app.models.Attendees.Delete(attendee.Id)

	if err != nil {
		app.serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendee deleted successfully"})
}
