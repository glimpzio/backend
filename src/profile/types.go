package profile

type NewUser struct {
	Id             *string  `json:"id,omitempty"`
	AuthId         string   `json:"authId"`
	Name           string   `json:"name"`
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
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Bio            string   `json:"bio"`
	ProfilePicture *string  `json:"profilePicture,omitempty"`
	Profile        *Profile `json:"profile"`
}
