package model

import "time"

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

type UserProfile struct {
	Id        string  `json:"id"`
	UserId    string  `json:"userId"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Website   *string `json:"website"`
	LinkedIn  *string `json:"linkedin"`
	Instagram *string `json:"instagram"`
	Facebook  *string `json:"facebook"`
}

type Link struct {
	Id            string    `json:"name"`
	UserProfileId string    `json:"userProfileId"`
	Url           string    `json:"url"`
	Expiry        time.Time `json:"expiry"`
}
