package profile

import "time"

type NewUser struct {
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	PersonalEmail  string   `json:"personalEmail"`
	Bio            string   `json:"bio"`
	ProfilePicture *string  `json:"profilePicture,omitempty"`
	Profile        *Profile `json:"profile"`
}

type Profile struct {
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Website  *string `json:"website,omitempty"`
	Linkedin *string `json:"linkedin,omitempty"`
}

type User struct {
	Id             string   `json:"id"`
	AuthId         string   `json:"authId"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Email          string   `json:"email"`
	Bio            string   `json:"bio"`
	ProfilePicture *string  `json:"profilePicture,omitempty"`
	Profile        *Profile `json:"profile"`
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
