package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (app *application) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		if authorizationHeader == "" {
			app.unauthorized(c)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			app.unauthorized(c)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if tokenString == authorizationHeader {
			app.unauthorized(c)
			c.Abort()
			return
		}
		id, err := app.verifyJWT(tokenString)
		if err != nil {
			app.unauthorized(c)
			c.Abort()
			return
		}

		user,err := app.models.Users.Get(*id)
		if err != nil {
			app.unauthorized(c)
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
