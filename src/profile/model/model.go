package model

import (
	"database/sql"
)

type Model struct {
	Db *sql.DB
}

// Create a new user
func (m *Model) CreateUser(name string, email string, bio string) (*User, error) {
	user := &User{}

	err := m.Db.QueryRow("INSERT INTO users (name, personalEmail, bio) VALUES ($1, $2, $3) RETURNING id, name, personalEmail, bio", name, email, bio).Scan(&user.Id, &user.Name, &user.PersonalEmail, &user.Bio)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user
func (m *Model) GetUser(id string) *User {
	user := &User{}

	err := m.Db.QueryRow("SELECT id, name, email, bio FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Email, &user.Bio)
	if err != nil {
		return nil
	}

	return user
}
