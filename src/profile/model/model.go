package model

import (
	"database/sql"
)

type Model struct {
	Db *sql.DB
}

// Create a new user
func (m *Model) CreateUser(id string, name string, personalEmail string, bio string, email *string, phone *string, website *string, linkedin *string) (*User, error) {
	user := &User{}

	tx, err := m.Db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow("SELECT id, auth_id, name, personal_email, bio, email, phone, website, linkedin FROM users WHERE auth_id = $1", id).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != sql.ErrNoRows {
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Commit()

		return user, nil
	}

	err = tx.QueryRow("INSERT INTO users (auth_id, name, personal_email, bio, email, phone, website, linkedin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, auth_id, name, personal_email, bio, email, phone, website, linkedin", id, name, personalEmail, bio, email, phone, website, linkedin).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
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

	err := m.Db.QueryRow("SELECT id, auth_id, name, personal_email, bio, email, phone, website, linkedin FROM users WHERE auth_id = $1", id).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil
	}

	return user
}
