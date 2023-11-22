package model

import (
	"database/sql"
	"time"
)

type Model struct {
	Db *sql.DB
}

// Create a new user
func (m *Model) CreateUser(authId string, firstName string, lastName string, personalEmail string, bio string, profilePicture *string, email *string, phone *string, website *string, linkedin *string) (*User, error) {
	user := &User{}

	tx, err := m.Db.Begin()
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

	err := m.Db.QueryRow("UPDATE users SET first_name = $1, last_name= $2, personal_email = $3, bio = $4, profile_picture = $5, email = $6, phone = $7, website = $8, linkedin = $9 WHERE id = $10 RETURNING id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin", firstName, lastName, personalEmail, bio, profilePicture, email, phone, website, linkedin, id).
		Scan(&user.Id, &user.AuthId, &user.FirstName, &user.LastName, &user.PersonalEmail, &user.Bio, &user.ProfilePicture, &user.Email, &user.Phone, &user.Website, &user.LinkedIn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by their id
func (m *Model) GetUserById(id string) (*User, error) {
	user := &User{}

	err := m.Db.QueryRow("SELECT id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE id = $1", id).
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

	err := m.Db.QueryRow("SELECT id, auth_id, first_name, last_name, personal_email, bio, profile_picture, email, phone, website, linkedin FROM users WHERE auth_id = $1", authId).
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

	err := m.Db.QueryRow("INSERT INTO invites (user_id, expires_at) VALUES ($1, $2) RETURNING id, user_id, expires_at", userId, expiresAt).
		Scan(&invite.Id, &invite.UserId, &invite.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// Get an invite
func (m *Model) GetInvite(id string) (*Invite, error) {
	invite := &Invite{}

	err := m.Db.QueryRow("SELECT id, user_id, expires_at FROM invites WHERE id = $1", id).
		Scan(&invite.Id, &invite.UserId, &invite.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return invite, nil
}

// Create a new email connection
func (m *Model) CreateEmailConnection(userId string, email string) (*EmailConnection, error) {
	emailConnection := &EmailConnection{}

	err := m.Db.QueryRow("INSERT INTO email_connections (user_id, email) VALUES ($1, $2) RETURNING id, user_id, email, connected_at", userId, email).
		Scan(&emailConnection.Id, &emailConnection.UserId, &emailConnection.Email, &emailConnection.ConnectedAt)
	if err != nil {
		return nil, err
	}

	return emailConnection, nil
}

// Get a list of distinct email connections for a user
func (m *Model) GetEmailConnections(userId string) ([]*EmailConnection, error) {
	emailConnections := []*EmailConnection{}

	rows, err := m.Db.Query("SELECT DISTINCT ON (email) id, user_id, email, connected_at FROM email_connections WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		emailConnection := &EmailConnection{}

		err := rows.Scan(&emailConnection.Id, &emailConnection.UserId, &emailConnection.Email, &emailConnection.ConnectedAt)
		if err != nil {
			return nil, err
		}

		emailConnections = append(emailConnections, emailConnection)
	}

	return emailConnections, nil
}
