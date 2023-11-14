// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Link struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	URL    string `json:"url"`
	Expiry string `json:"expiry"`
}

type NewLink struct {
	UserID string `json:"userId"`
}

type NewUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

type User struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Bio         string       `json:"bio"`
	UserContact *UserContact `json:"userContact"`
}

type UserContact struct {
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Website  *string `json:"website,omitempty"`
	Linkedin *string `json:"linkedin,omitempty"`
}
