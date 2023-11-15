package profile

type NewUser struct {
	Name          string  `json:"name"`
	PersonalEmail string  `json:"personalEmail"`
	Bio           string  `json:"bio"`
	Profile       Profile `json:"profile,omitempty"`
}

type Profile struct {
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Website  *string `json:"website,omitempty"`
	Linkedin *string `json:"linkedin,omitempty"`
}

type User struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Bio     string   `json:"bio"`
	Profile *Profile `json:"profile"`
}
