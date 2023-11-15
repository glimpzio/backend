package model

import (
	"database/sql"
)

type Model struct {
	Db *sql.DB
}

// Create a new user
func (m *Model) CreateUser(authId string, name string, email string, bio string) (*User, error) {
	user := &User{}

	tx, err := m.Db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow("SELECT * FROM users WHERE auth_id = $1", authId).Scan(&user.Id, &user.Name, &user.PersonalEmail, &user.Bio)
	if err != sql.ErrNoRows {
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Commit()

		return user, nil
	}

	err = tx.QueryRow("INSERT INTO users (auth_id, name, personal_email, bio) VALUES ($1, $2, $3, $4) RETURNING id, auth_id, name, personal_email, bio", authId, name, email, bio).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by their id
func (m *Model) GetUser(id string) *User {
	user := &User{}

	err := m.Db.QueryRow("SELECT id, auth_id, name, email, bio FROM users WHERE id = $1", id).Scan(&user.Id, &user.AuthId, &user.Name, &user.Email, &user.Bio)
	if err != nil {
		return nil
	}

	return user
}
