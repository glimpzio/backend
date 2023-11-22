package model

import (
	"database/sql"
)

type Model struct {
	Db *sql.DB
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
