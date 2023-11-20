package model

import "time"

type User struct {
	Id             string  `json:"id"`
	AuthId         string  `json:"authId"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	PersonalEmail  string  `json:"personalEmail"`
	Bio            string  `json:"bio"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Website        *string `json:"website,omitempty"`
	LinkedIn       *string `json:"linkedin,omitempty"`
}

type Invite struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type EmailConnection struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Email       string    `json:"email"`
	ConnectedAt time.Time `json:"connectedAt"`
}
