package model

import (
	"database/sql"
	"time"
)

type Model struct {
	Db *sql.DB
}

// Create a new user
func (m *Model) CreateUser(authId string, name string, personalEmail string, bio string, profilePicture *string, email *string, phone *string, website *string, linkedin *string) (*User, error) {
	user := &User{}

	tx, err := m.Db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow("SELECT id, auth_id, name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE auth_id = $1", authId).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != sql.ErrNoRows {
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Commit()

		return user, nil
	}

	err = tx.QueryRow("INSERT INTO users (auth_id, name, personal_email, bio, profile_picture, email, phone, website, linkedin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, auth_id, name, personal_email, bio, profile_picture, email, phone, website, linkedin", authId, name, personalEmail, bio, profilePicture, email, phone, website, linkedin).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update a user
func (m *Model) UpdateUser(id string, name string, personalEmail string, bio string, profilePicture *string, email *string, phone *string, website *string, linkedin *string) (*User, error) {
	user := &User{}

	err := m.Db.QueryRow("UPDATE users SET name = $1, personal_email = $2, bio = $3, profile_picture = $4, email = $5, phone = $6, website = $7, linkedin = $8 WHERE id = $9 RETURNING id, auth_id, name, personal_email, bio, profile_picture, email, phone, website, linkedin", name, personalEmail, bio, profilePicture, email, phone, website, linkedin, id).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by their id
func (m *Model) GetUserById(id string) (*User, error) {
	user := &User{}

	err := m.Db.QueryRow("SELECT id, auth_id, name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE id = $1", id).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by their auth id
func (m *Model) GetUserByAuthId(authId string) (*User, error) {
	user := &User{}

	err := m.Db.QueryRow("SELECT id, auth_id, name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE auth_id = $1", authId).Scan(&user.Id, &user.AuthId, &user.Name, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create a new link
func (m *Model) CreateLink(userId string, expiresAt time.Time) (*Link, error) {
	link := &Link{}

	err := m.Db.QueryRow("INSERT INTO links (user_id, expires_at) VALUES ($1, $2) RETURNING id, user_id, expires_at", userId, expiresAt).Scan(&link.Id, &link.UserId, &link.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// Get a new link
func (m *Model) GetLink(id string) (*Link, error) {
	link := &Link{}

	err := m.Db.QueryRow("SELECT id, user_id, expires_at FROM links WHERE id = $1", id).Scan(&link.Id, &link.UserId, &link.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return link, nil
}
