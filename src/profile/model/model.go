package model

import (
	"database/sql"
)

type Model struct {
	Db *sql.DB
}

// Create a new user
func (m *Model) CreateUser(name string, email string, bio string) (*User, *UserProfile, error) {
	user := &User{}
	profile := &UserProfile{}

	tx, err := m.Db.Begin()
	if err != nil {
		return nil, nil, err
	}

	err = tx.QueryRow("INSERT INTO users (name, email, bio) VALUES ($1, $2, $3) RETURNING id, name, email, bio", name, email, bio).Scan(&user.Id, &user.Name, &user.Email, &user.Bio)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.QueryRow("INSERT INTO user_profiles (user_id) VALUES ($1) RETURNING id, user_id, email, phone, website, linkedin", user.Id).Scan(&profile.Id, &profile.UserId, &profile.Email, &profile.Phone, &profile.Website, &profile.LinkedIn)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	return user, profile, nil
}

// Get a user
func (m *Model) GetUser(id string) (*User, *UserProfile) {
	user := &User{}
	profile := &UserProfile{}

	err := m.Db.QueryRow("SELECT id, name, email, bio FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Email, &user.Bio)
	if err != nil {
		return nil, nil
	}

	err = m.Db.QueryRow("SELECT id, user_id, email, phone, website, linkedin FROM user_profiles WHERE user_id = $1", user.Id).Scan(&profile.Id, &profile.UserId, &profile.Email, &profile.Phone, &profile.Website, &profile.LinkedIn)
	if err != nil {
		return nil, nil
	}

	return user, profile
}
