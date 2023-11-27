package model

import (
	"database/sql"
)

type Model struct {
	ReadDb  *sql.DB
	WriteDb *sql.DB
}

// Create a new custom connection
func (m *Model) CreateCustomConnection(userId string, firstName *string, lastName *string, notes *string, email *string, phone *string, website *string, linkedin *string) (*CustomConnection, error) {
	customConnection := &CustomConnection{}

	err := m.WriteDb.QueryRow("INSERT INTO custom_connections (user_id, first_name, last_name, notes, email, phone, website, linkedin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, user_id, connected_at, first_name, last_name, notes, email, phone, website, linkedin",
		userId, firstName, lastName, notes, email, phone, website, linkedin).
		Scan(&customConnection.Id, &customConnection.UserId, &customConnection.ConnectedAt, &customConnection.FirstName, &customConnection.LastName, &customConnection.Notes, &customConnection.Email, &customConnection.Phone, &customConnection.Website, &customConnection.LinkedIn)
	if err != nil {
		return nil, err
	}

	return customConnection, nil
}

// Update a new custom connection
func (m *Model) UpdateCustomConnection(id string, firstName *string, lastName *string, notes *string, email *string, phone *string, website *string, linkedin *string) (*CustomConnection, error) {
	customConnection := &CustomConnection{}

	err := m.WriteDb.QueryRow("UPDATE custom_connections SET first_name = $1, last_name = $2, notes = $3, email = $4, phone = $5, website = $6, linkedin = $7 WHERE id = $8 RETURNING id, user_id, connected_at, first_name, last_name, notes, email, phone, website, linkedin",
		firstName, lastName, notes, email, phone, website, linkedin, id).
		Scan(&customConnection.Id, &customConnection.UserId, &customConnection.ConnectedAt, &customConnection.FirstName, &customConnection.LastName, &customConnection.Notes, &customConnection.Email, &customConnection.Phone, &customConnection.Website, &customConnection.LinkedIn)
	if err != nil {
		return nil, err
	}

	return customConnection, nil
}

// Delete a custom connection
func (m *Model) DeleteCustomConnection(id string) (*CustomConnection, error) {
	customConnection := &CustomConnection{}

	err := m.WriteDb.QueryRow("DELETE FROM custom_connections WHERE id = $1 RETURNING id, user_id, connected_at, first_name, last_name, notes, email, phone, website, linkedin", id).
		Scan(&customConnection.Id, &customConnection.UserId, &customConnection.ConnectedAt, &customConnection.FirstName, &customConnection.LastName, &customConnection.Notes, &customConnection.Email, &customConnection.Phone, &customConnection.Website, &customConnection.LinkedIn)
	if err != nil {
		return nil, err
	}

	return customConnection, nil
}

// Get a custom connection
func (m *Model) GetCustomConnection(id string) (*CustomConnection, error) {
	customConnection := &CustomConnection{}

	err := m.ReadDb.QueryRow("SELECT id, user_id, connected_at, first_name, last_name, notes, email, phone, website, linkedin FROM custom_connections WHERE id = $1", id).
		Scan(&customConnection.Id, &customConnection.UserId, &customConnection.ConnectedAt, &customConnection.FirstName, &customConnection.LastName, &customConnection.Notes, &customConnection.Email, &customConnection.Phone, &customConnection.Website, &customConnection.LinkedIn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return customConnection, nil
}

// Get a list of custom connections for a user
func (m *Model) GetCustomConnections(userId string) ([]*CustomConnection, error) {
	customConnections := []*CustomConnection{}

	rows, err := m.ReadDb.Query("SELECT id, user_id, connected_at, first_name, last_name, notes, email, phone, website, linkedin FROM custom_connections WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		customConnection := &CustomConnection{}

		err := rows.Scan(&customConnection.Id, &customConnection.UserId, &customConnection.ConnectedAt, &customConnection.FirstName, &customConnection.LastName, &customConnection.Notes, &customConnection.Email, &customConnection.Phone, &customConnection.Website, &customConnection.LinkedIn)
		if err != nil {
			return nil, err
		}

		customConnections = append(customConnections, customConnection)
	}

	return customConnections, nil
}
