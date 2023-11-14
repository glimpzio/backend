package profile

import (
	"database/sql"

	"github.com/glimpzio/backend/profile/model"
)

type ProfileService struct {
	model *model.Model
}

// Create a new profile service
func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{model: &model.Model{Db: db}}
}

// Create a new user
func (p *ProfileService) NewUser(user *NewUser) (*User, error) {
	rawUser, rawProfile, err := p.model.CreateUser(user.Name, user.Email, user.Bio)

	if err != nil {
		return nil, err
	}

	return &User{Id: rawUser.Id, Name: rawUser.Name, Email: rawUser.Email, Bio: rawUser.Bio, Profile: &Profile{Email: rawProfile.Email, Phone: rawProfile.Phone, Website: rawProfile.Website, Linkedin: rawProfile.LinkedIn}}, nil
}

// Get a user
func (p *ProfileService) GetUser(userId string) *User {
	rawUser, rawProfile := p.model.GetUser(userId)

	return &User{Id: rawUser.Id, Name: rawUser.Name, Email: rawUser.Email, Bio: rawUser.Bio, Profile: &Profile{Email: rawProfile.Email, Phone: rawProfile.Phone, Website: rawProfile.Website, Linkedin: rawProfile.LinkedIn}}
}
