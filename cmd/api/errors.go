package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) errorResponse(c *gin.Context, status int, message string, details string) {
	c.JSON(status, gin.H{
		"error":   message,
		"details": details,
	})
}

func (app *application) badRequest(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, "bad request", err.Error())
}

func (app *application) serverError(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusInternalServerError, "internal server error", err.Error())
}

func (app *application) notFound(c *gin.Context) {
	app.errorResponse(c, http.StatusNotFound, "resource not found", "")
}

func (app *application) unauthorized(c *gin.Context) {
	app.errorResponse(c, http.StatusUnauthorized, "unauthorized", "")
}

func (app *application) invalidCredentials(c *gin.Context) {
	app.errorResponse(c, http.StatusUnauthorized, "invalid credentials", "")
}

func (app *application) editConflict(c *gin.Context) {
	app.errorResponse(c, http.StatusConflict, "edit conflict", "")
}