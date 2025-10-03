package main

import (
	"first-rest-api/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createEvent godoc
// @Summary      Create an event
// @Description  Create a new event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event body database.Event true "Event to create"
// @Success      201  {object} database.Event
// @Failure      400  {object} ErrorResponse
// @Failure      401  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /api/v1/events [post]
func (app *application) createEvent(ctx *gin.Context) {
	var event database.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		app.badRequest(ctx, err)
		return
	}

	user, err := app.getUser(ctx)

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	event.OwnerId = user.Id

	err = app.models.Events.Insert(&event)

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

// getAllEvents godoc
// @Summary      Get all events
// @Description  Get a list of all events
// @Tags         Events
// @Accept       json
// @Produce      json
// @Success      200  {array}  database.Event
// @Failure      500  {object} ErrorResponse
// @Router       /api/v1/events [get]
func (app *application) getAllEvents(ctx *gin.Context) {

	events, err := app.models.Events.GetAll()

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, events)
}

// getEvent godoc
// @Summary      Get event
// @Description  Get a single event by ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id  path     int     true  "Event ID"
// @Success      200  {object}  database.Event
// @Failure      400  {object} ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/events/{id} [get]
func (app *application) getEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	event, err := app.models.Events.Get(id)

	if err != nil {
		app.notFound(ctx)
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// updateEvent godoc
// @Summary      Update an event
// @Description  Update an existing event by ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     int             true  "Event ID"
// @Param        event body     database.Event true  "Updated event data"
// @Success      200   {object} database.Event
// @Failure      400   {object} ErrorResponse
// @Failure      401   {object} ErrorResponse
// @Failure      404   {object} ErrorResponse
// @Failure      409   {object} ErrorResponse
// @Failure      500   {object} ErrorResponse
// @Router       /api/v1/events/{id} [put]
func (app *application) updateEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		app.notFound(ctx)
		return
	}

	user, err := app.getUser(ctx)

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	if existingEvent.OwnerId != user.Id {
		app.unauthorized(ctx)
		return
	}

	var updatedEvent database.Event

	err = app.readJSON(ctx.Writer, ctx.Request, &updatedEvent)
	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(&updatedEvent); err != nil {
		app.editConflict(ctx)
		return
	}

	ctx.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent godoc
// @Summary      Delete an event
// @Description  Delete an existing event by ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     int     true  "Event ID"
// @Success      200   {object}  gin.H
// @Failure      400   {object} ErrorResponse
// @Failure      401   {object} ErrorResponse
// @Failure      404   {object} ErrorResponse
// @Failure      500   {object} ErrorResponse
// @Router       /events/{id} [delete]
func (app *application) deleteEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		app.notFound(ctx)
		return
	}

	user, err := app.getUser(ctx)

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	if existingEvent.OwnerId != user.Id {
		app.unauthorized(ctx)
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})

}
