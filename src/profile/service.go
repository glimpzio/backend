package profile

import (
	"database/sql"
	"errors"
	"time"

	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile/model"
)

type ProfileService struct {
	model    *model.Model
	mailList *misc.MailList
}

// Create a new profile service
func NewProfileService(db *sql.DB, mailList *misc.MailList) *ProfileService {
	return &ProfileService{model: &model.Model{Db: db}, mailList: mailList}
}

// Upsert a user
func (p *ProfileService) UpsertUser(authId string, user *NewUser) (*User, error) {
	var rawUser *model.User
	var err error

	existing, err := p.model.GetUserByAuthId(authId)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		rawUser, err = p.model.CreateUser(authId, user.FirstName, user.LastName, user.PersonalEmail, user.Bio, user.ProfilePicture, user.Profile.Email, user.Profile.Phone, user.Profile.Website, user.Profile.Linkedin)
		if err != nil {
			return nil, err
		}

		err = p.mailList.AddAccount(user.PersonalEmail, user.FirstName, user.LastName)
	} else {
		rawUser, err = p.model.UpdateUser(existing.Id, user.FirstName, user.LastName, user.PersonalEmail, user.Bio, user.ProfilePicture, user.Profile.Email, user.Profile.Phone, user.Profile.Website, user.Profile.Linkedin)
	}

	if err != nil {
		return nil, err
	}

	return &User{
		Id:             rawUser.Id,
		AuthId:         rawUser.AuthId,
		FirstName:      rawUser.FirstName,
		LastName:       rawUser.LastName,
		Email:          rawUser.PersonalEmail,
		Bio:            rawUser.Bio,
		ProfilePicture: user.ProfilePicture,
		Profile: &Profile{
			Email:    rawUser.Email,
			Phone:    rawUser.Phone,
			Website:  rawUser.Website,
			Linkedin: rawUser.LinkedIn,
		},
	}, nil
}

// Get a user by id
func (p *ProfileService) GetUserById(id string) (*User, error) {
	rawUser, err := p.model.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return &User{
		Id:             rawUser.Id,
		AuthId:         rawUser.AuthId,
		FirstName:      rawUser.FirstName,
		LastName:       rawUser.LastName,
		Email:          rawUser.PersonalEmail,
		Bio:            rawUser.Bio,
		ProfilePicture: rawUser.ProfilePicture,
		Profile: &Profile{
			Email:    rawUser.Email,
			Phone:    rawUser.Phone,
			Website:  rawUser.Website,
			Linkedin: rawUser.LinkedIn,
		},
	}, nil
}

// Get a user by auth id
func (p *ProfileService) GetUserByAuthId(authId string) (*User, error) {
	rawUser, err := p.model.GetUserByAuthId(authId)

	if err != nil {
		return nil, err
	}

	return &User{
		Id:             rawUser.Id,
		AuthId:         rawUser.AuthId,
		FirstName:      rawUser.FirstName,
		LastName:       rawUser.LastName,
		Email:          rawUser.PersonalEmail,
		Bio:            rawUser.Bio,
		ProfilePicture: rawUser.ProfilePicture,
		Profile: &Profile{
			Email:    rawUser.Email,
			Phone:    rawUser.Phone,
			Website:  rawUser.Website,
			Linkedin: rawUser.LinkedIn,
		},
	}, nil
}

// Create a new invite
func (p *ProfileService) CreateInvite(userId string) (*Invite, error) {
	expiresAt := time.Time.Add(time.Now(), time.Duration(time.Duration.Hours(24*3)))

	rawInvite, err := p.model.CreateInvite(userId, expiresAt)
	if err != nil {
		return nil, err
	}

	return &Invite{
		Id:        rawInvite.Id,
		UserId:    rawInvite.UserId,
		ExpiresAt: rawInvite.ExpiresAt,
	}, nil
}

// Get an invite
func (p *ProfileService) GetInvite(id string) (*Invite, *User, error) {
	rawInvite, err := p.model.GetInvite(id)
	if err != nil {
		return nil, nil, err
	}

	if rawInvite.ExpiresAt.Compare(time.Now()) < 0 {
		return nil, nil, errors.New("invite has expired")
	}

	rawUser, err := p.model.GetUserById(rawInvite.UserId)
	if err != nil {
		return nil, nil, err
	}

	return &Invite{
			Id:        rawInvite.Id,
			UserId:    rawInvite.UserId,
			ExpiresAt: rawInvite.ExpiresAt,
		}, &User{
			Id:             rawUser.Id,
			AuthId:         rawUser.AuthId,
			FirstName:      rawUser.FirstName,
			LastName:       rawUser.LastName,
			Email:          rawUser.PersonalEmail,
			Bio:            rawUser.Bio,
			ProfilePicture: rawUser.ProfilePicture,
			Profile: &Profile{
				Email:    rawUser.Email,
				Phone:    rawUser.Phone,
				Website:  rawUser.Website,
				Linkedin: rawUser.LinkedIn,
			},
		}, nil
}
