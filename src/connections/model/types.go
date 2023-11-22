package model

import "time"

type EmailConnection struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Email       string    `json:"email"`
	ConnectedAt time.Time `json:"connectedAt"`
}
