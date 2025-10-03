package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.Id)

}

func (m *UserModel) get(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, args...)
	var user User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {

	query := "SELECT id, name, email, password FROM users WHERE email = $1"

	return m.get(query, email)
}

func (m *UserModel) Get(id int) (*User, error) {

	query := "SELECT id, name, email, password FROM users WHERE id = $1"

	user, err := m.get(query, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}
