package model

type User struct {
	Id             string  `json:"id"`
	AuthId         string  `json:"authId"`
	Name           string  `json:"name"`
	PersonalEmail  string  `json:"personalEmail"`
	Bio            string  `json:"bio"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Website        *string `json:"website,omitempty"`
	LinkedIn       *string `json:"linkedin,omitempty"`
}
