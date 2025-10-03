package main

import (
	"first-rest-api/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) loginUser(c *gin.Context) {
	var login loginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		app.badRequest(c, err)
		return
	}

	user, err := app.models.Users.GetByEmail(login.Email)

	if err != nil {
		app.badRequest(c, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))

	if err != nil {
		app.badRequest(c, err)
		return
	}

	token, err := app.generateJWT(user.Id)

	if err != nil {
		app.serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token})
}

func (m *application) registerUser(c *gin.Context) {
	var register registerRequest

	err := m.readJSON(c.Writer, c.Request, &register)
	if err != nil {
		m.badRequest(c, err)
		return
	}

	log.Print(register.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)

	if err != nil {
		m.serverError(c, err)
		return
	}

	register.Password = string(hashedPassword)

	user := database.User{
		Name:     register.Name,
		Email:    register.Email,
		Password: register.Password,
	}

	err = m.models.Users.Insert(&user)

	if err != nil {
		log.Print(err.Error())
		m.serverError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}
