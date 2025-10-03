package main

import (
	"first-rest-api/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(ctx *gin.Context) {
	var event database.Event

	err := app.readJSON(ctx.Writer, ctx.Request, &event)
	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	err = app.models.Events.Insert(&event)

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func (app *application) getAllEvents(ctx *gin.Context) {

	events, err := app.models.Events.GetAll()

	if err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, events)
}

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

func (app *application) updateEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	_, err = app.models.Events.Get(id)

	if err != nil {
		app.notFound(ctx)
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

func (app *application) deleteEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		app.badRequest(ctx, err)
		return
	}

	_, err = app.models.Events.Get(id)

	if err != nil {
		app.notFound(ctx)
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})

}
