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

// Upsert a user
func (p *ProfileService) UpsertUser(user *NewUser) (*User, error) {
	var rawUser *model.User
	var err error

	if user.Id == nil {
		rawUser, err = p.model.CreateUser(user.AuthId, user.Name, user.PersonalEmail, user.Bio, user.ProfilePicture, user.Profile.Email, user.Profile.Phone, user.Profile.Website, user.Profile.Linkedin)
	} else {
		rawUser, err = p.model.UpdateUser(*user.Id, user.Name, user.PersonalEmail, user.Bio, user.ProfilePicture, user.Profile.Email, user.Profile.Phone, user.Profile.Website, user.Profile.Linkedin)
	}

	if err != nil {
		return nil, err
	}

	return &User{Id: rawUser.Id, AuthId: rawUser.AuthId, Name: rawUser.Name, Email: rawUser.PersonalEmail, Bio: rawUser.Bio, ProfilePicture: user.ProfilePicture, Profile: &Profile{Email: rawUser.Email, Phone: rawUser.Phone, Website: rawUser.Website, Linkedin: rawUser.LinkedIn}}, nil
}

// Get a user by id
func (p *ProfileService) GetUserById(id string) (*User, error) {
	rawUser, err := p.model.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return &User{Id: rawUser.Id, AuthId: rawUser.AuthId, Name: rawUser.Name, Email: rawUser.PersonalEmail, Bio: rawUser.Bio, ProfilePicture: rawUser.ProfilePicture, Profile: &Profile{Email: rawUser.Email, Phone: rawUser.Phone, Website: rawUser.Website, Linkedin: rawUser.LinkedIn}}, nil
}

// Get a user by auth id
func (p *ProfileService) GetUserByAuthId(authId string) (*User, error) {
	rawUser, err := p.model.GetUserByAuthId(authId)

	if err != nil {
		return nil, err
	}

	return &User{Id: rawUser.Id, AuthId: rawUser.AuthId, Name: rawUser.Name, Email: rawUser.PersonalEmail, Bio: rawUser.Bio, ProfilePicture: rawUser.ProfilePicture, Profile: &Profile{Email: rawUser.Email, Phone: rawUser.Phone, Website: rawUser.Website, Linkedin: rawUser.LinkedIn}}, nil
}
