package main

import (
	"first-rest-api/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// addAttendee godoc
// @Summary      Add an attendee
// @Description  Add an attendee to an event
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     int     true  "Event ID"
// @Param        userId    path     int     true  "User ID"
// @Success      200  {object}  gin.H
// @Failure      400  {object} ErrorResponse
// @Failure      401  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /api/v1/events/{id}/attendees/{userId} [post]

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

	event, err := app.models.Events.Get(eventID)
	if err != nil {
		app.serverError(c, err)
		return
	}

	if event == nil {
		app.notFound(c)
		return
	}

	user, err := app.getUser(c)

	if err != nil {
		app.serverError(c, err)
		return
	}

	if user.Id != event.OwnerId {
		app.unauthorized(c)
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

// getAttendees godoc
// @Summary      Get attendees
// @Description  Get a list of attendees for an event
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Success      200  {array}  []database.User
// @Failure      500  {object} ErrorResponse
// @Router       /api/v1/events/{id}/attendees [get]
func (m *application) getAttendees(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		m.badRequest(c, err)
		return
	}

	user, err := m.getUser(c)

	if err != nil {
		m.serverError(c, err)
		return
	}

	even, err := m.models.Events.Get(eventID)

	if err != nil {
		m.serverError(c, err)
		return
	}

	if even.OwnerId != user.Id {
		m.unauthorized(c)
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

// deleteAttendee godoc
// @Summary      Delete an attendee
// @Description  Delete an attendee from an event
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     int     true  "Event ID"
// @Param        userId    path     int     true  "User ID"
// @Success      200  {object}  gin.H
// @Failure      400  {object} ErrorResponse
// @Failure      401  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /api/v1/events/{id}/attendees/{userId} [delete]
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

	event, err := app.models.Events.Get(eventID)
	if err != nil {
		app.serverError(c, err)
		return
	}

	if event == nil {
		app.notFound(c)
		return
	}

	user, err := app.getUser(c)

	if err != nil {
		app.serverError(c, err)
		return
	}

	if user.Id != event.OwnerId {
		app.unauthorized(c)
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
