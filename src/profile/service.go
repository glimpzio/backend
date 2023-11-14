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
	rawUser, err := p.model.CreateUser(user.Name, user.PersonalEmail, user.Bio)

	if err != nil {
		return nil, err
	}

	return &User{Id: rawUser.Id, Name: rawUser.Name, Email: rawUser.PersonalEmail, Bio: rawUser.Bio, Profile: &Profile{Email: rawUser.Email, Phone: rawUser.Phone, Website: rawUser.Website, Linkedin: rawUser.LinkedIn}}, nil
}

// Get a user
func (p *ProfileService) GetUser(userId string) *User {
	rawUser := p.model.GetUser(userId)

	if rawUser == nil {
		return nil
	}

	return &User{Id: rawUser.Id, Name: rawUser.Name, Email: rawUser.PersonalEmail, Bio: rawUser.Bio, Profile: &Profile{Email: rawUser.Email, Phone: rawUser.Phone, Website: rawUser.Website, Linkedin: rawUser.LinkedIn}}
}
