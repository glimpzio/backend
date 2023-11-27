package model

import "time"

type CustomConnection struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	ConnectedAt time.Time `json:"connectedAt"`
	FirstName   *string   `json:"first_name,omitempty"`
	LastName    *string   `json:"last_name,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	Website     *string   `json:"website,omitempty"`
	LinkedIn    *string   `json:"linkedin,omitempty"`
}
