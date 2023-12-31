package profile

import (
	"database/sql"
	"time"

	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile/model"
)

type ProfileService struct {
	model    *model.Model
	mailList *misc.MailList
}

// Create a new profile service
func NewProfileService(readDb *sql.DB, writeDb *sql.DB, mailList *misc.MailList) *ProfileService {
	return &ProfileService{model: &model.Model{ReadDb: readDb, WriteDb: writeDb}, mailList: mailList}
}

// Upsert a user
func (p *ProfileService) UpsertUser(authId string, newUser *NewUser) (*User, error) {
	var rawUser *model.User

	existing, err := p.model.GetUserByAuthId(authId)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		user, err := p.model.CreateUser(authId, newUser.FirstName, newUser.LastName, newUser.PersonalEmail, newUser.Bio, newUser.ProfilePicture, newUser.Profile.Email, newUser.Profile.Phone, newUser.Profile.Website, newUser.Profile.LinkedIn)
		if err != nil {
			return nil, err
		}

		rawUser = user

		if err := p.mailList.AddAccount(newUser.PersonalEmail, newUser.FirstName, newUser.LastName); err != nil {
			return nil, err
		}
	} else {
		user, err := p.model.UpdateUser(existing.Id, newUser.FirstName, newUser.LastName, newUser.PersonalEmail, newUser.Bio, newUser.ProfilePicture, newUser.Profile.Email, newUser.Profile.Phone, newUser.Profile.Website, newUser.Profile.LinkedIn)
		if err != nil {
			return nil, err
		}

		rawUser = user
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
			LinkedIn: rawUser.LinkedIn,
		},
	}, nil
}

// Get a user by id
func (p *ProfileService) GetUserById(id string) (*User, error) {
	rawUser, err := p.model.GetUserById(id)
	if err != nil {
		return nil, err
	} else if rawUser == nil {
		return nil, ErrInvalidUser
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
			LinkedIn: rawUser.LinkedIn,
		},
	}, nil
}

// Get a user by auth id
func (p *ProfileService) GetUserByAuthId(authId string) (*User, error) {
	rawUser, err := p.model.GetUserByAuthId(authId)
	if err != nil {
		return nil, err
	} else if rawUser == nil {
		return nil, ErrInvalidUser
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
			LinkedIn: rawUser.LinkedIn,
		},
	}, nil
}

// Create a new invite
func (p *ProfileService) CreateInvite(userId string) (*Invite, error) {
	currentTime := time.Now()
	expiresAt := currentTime.Add(24 * time.Hour)

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
	} else if rawInvite == nil {
		return nil, nil, ErrInvalidInvite
	}

	if rawInvite.ExpiresAt.Compare(time.Now()) < 0 {
		return nil, nil, ErrInviteExpired
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
				LinkedIn: rawUser.LinkedIn,
			},
		}, nil
}
