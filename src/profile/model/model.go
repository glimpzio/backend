package model

import (
	"database/sql"
	"time"

	"github.com/glimpzio/backend/misc"
)

type Model struct {
	ReadDb  *sql.DB
	WriteDb *sql.DB
}

// Create a new user
func (m *Model) CreateUser(authId string, firstName string, lastName string, personalEmail string, bio string, profilePicture *string, email *string, phone *string, website *string, linkedin *string) (*User, error) {
	user := &User{}

	tx, err := m.WriteDb.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow("SELECT id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE auth_id = $1", authId).
		Scan(&user.Id, &user.AuthId, &user.FirstName, &user.LastName, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != sql.ErrNoRows {
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Commit()

		return user, nil
	}

	err = tx.QueryRow("INSERT INTO users (auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin", authId, firstName, lastName, personalEmail, bio, profilePicture, email, phone, website, linkedin).
		Scan(&user.Id, &user.AuthId, &user.FirstName, &user.LastName, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
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
func (m *Model) UpdateUser(id string, firstName string, lastName string, personalEmail string, bio string, profilePicture *string, email *string, phone *string, website *string, linkedin *string) (*User, error) {
	user := &User{}

	err := m.WriteDb.QueryRow("UPDATE users SET first_name = $1, last_name= $2, personal_email = $3, bio = $4, profile_picture = $5, email = $6, phone = $7, website = $8, linkedin = $9 WHERE id = $10 RETURNING id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin", firstName, lastName, personalEmail, bio, profilePicture, email, phone, website, linkedin, id).
		Scan(&user.Id, &user.AuthId, &user.FirstName, &user.LastName, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by their id
func (m *Model) GetUserById(id string) (*User, error) {
	user := &User{}

	err := m.ReadDb.QueryRow("SELECT id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE id = $1", id).
		Scan(&user.Id, &user.AuthId, &user.FirstName, &user.LastName, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

// Get a user by their auth id
func (m *Model) GetUserByAuthId(authId string) (*User, error) {
	user := &User{}

	err := m.ReadDb.QueryRow("SELECT id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE auth_id = $1", authId).
		Scan(&user.Id, &user.AuthId, &user.FirstName, &user.LastName, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

// Create a new invite
func (m *Model) CreateInvite(userId string, expiresAt time.Time) (*Invite, error) {
	invite := &Invite{}

	err := m.WriteDb.QueryRow("INSERT INTO invites (user_id, expires_at) VALUES ($1, $2) RETURNING id, user_id, expires_at", userId, misc.FormatTime(expiresAt)).
		Scan(&invite.Id, &invite.UserId, &invite.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// Get an invite
func (m *Model) GetInvite(id string) (*Invite, error) {
	invite := &Invite{}

	err := m.ReadDb.QueryRow("SELECT id, user_id, expires_at FROM invites WHERE id = $1", id).
		Scan(&invite.Id, &invite.UserId, &invite.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return invite, nil
}
