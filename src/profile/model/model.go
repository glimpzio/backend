package model

import "time"

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

type UserContact struct {
	UserId   string  `json:"userId"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Website  *string `json:"website"`
	LinkedIn *string `json:"linkedin"`
}

type Link struct {
	Id     string    `json:"name"`
	UserId string    `json:"userId"`
	Url    string    `json:"url"`
	Expiry time.Time `json:"expiry"`
}
