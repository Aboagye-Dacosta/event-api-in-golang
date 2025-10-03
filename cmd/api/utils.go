package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

// readJSON reads and decodes a JSON request body into the provided data structure.
// It ensures that the request body does not exceed 1MB and that it contains only
// a single JSON value.
//
// Parameters:
//   - w: The http.ResponseWriter to send error responses if necessary.
//   - r: The http.Request containing the JSON body to be read.
//   - data: A pointer to the data structure where the JSON will be unmarshaled.
//
// Returns:
//   - An error if the JSON decoding fails, the body exceeds the size limit, or
//     if there is more than one JSON value in the body.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

func (app *application) generateJWT(userId int) (string, error) {
	// Define the JWT claims
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(app.jwtSecrete))
}

func (app *application) verifyJWT(tokenString string) (*int, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(app.jwtSecrete), nil
	})

	if err != nil || !token.Valid {
		if !token.Valid {
			return nil, errors.New("invalid token")
		}

		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userId, ok := claims["userId"].(float64)

	if !ok {
		return nil, errors.New("could not parse userId")
	}
	userIdInt := int(userId)
	return &userIdInt, nil
}
