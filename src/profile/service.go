package profile

import (
	"database/sql"
	"fmt"
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
	} else if rawUser == nil {
		return nil, ErrDoesNotExist
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
	} else if rawUser == nil {
		return nil, ErrDoesNotExist
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
	expiresAt := time.Time.Add(time.Now(), time.Duration(time.Duration.Hours(24)))

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
		return nil, nil, ErrDoesNotExist
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
				Linkedin: rawUser.LinkedIn,
			},
		}, nil
}

// Connec the users by email signup
func (p *ProfileService) ConnectByEmail(inviteId string, email string, subscribe bool) (*EmailConnection, error) {
	rawInvite, err := p.model.GetInvite(inviteId)
	if err != nil {
		return nil, err
	} else if rawInvite == nil {
		return nil, ErrDoesNotExist
	}

	if rawInvite.ExpiresAt.Compare(time.Now()) < 0 {
		return nil, ErrInviteExpired
	}

	user, err := p.model.GetUserById(rawInvite.UserId)
	if err != nil {
		return nil, err
	}

	rawEmailConnection, err := p.model.CreateEmailConnection(user.Id, email)
	if err != nil {
		return nil, err
	}

	if subscribe {
		err = p.mailList.AddMarketing(email, nil, nil)
		if err != nil {
			return nil, err
		}
	}

	body := fmt.Sprintf("Hey, hope you're well!\n\nAs you requested, here's the Glimpz profile for %s %s:\n\n- Bio: %s", user.FirstName, user.LastName, user.Bio)

	if user.Email != nil {
		body += fmt.Sprintf("\n- Email: %s", *user.Email)
	}
	if user.Phone != nil {
		body += fmt.Sprintf("\n- Phone: %s", *user.Phone)
	}
	if user.Website != nil {
		body += fmt.Sprintf("\n- Website: %s", *user.Website)
	}
	if user.LinkedIn != nil {
		body += fmt.Sprintf("\n- LinkedIn: %s", *user.LinkedIn)
	}

	body += fmt.Sprintf("\n\nWe've also forwarded your email to %s so they can follow up with you when they get a chance.", user.FirstName)
	body += fmt.Sprintf("\n\nBy the way, if you're ever looking to increase your own leads and sales from networking events like %s, did you know that you can make your own Glimpz profile for free right now? Glimpz makes it easy for you to connect with other professionals at networking events and convert them into long-lasting business partners or clients. Check it out at https://glimpz.io?referral=uide-%s", user.FirstName, user.Id)

	body += "\n\nWarm regards,\nBen"

	subject := fmt.Sprintf("Here's %s %s's Profile As You Requested!", user.FirstName, user.LastName)

	err = p.mailList.SendMail(email, email, subject, body)
	if err != nil {
		return nil, err
	}

	return &EmailConnection{
		Id:          rawEmailConnection.Id,
		UserId:      rawEmailConnection.UserId,
		Email:       rawEmailConnection.Email,
		ConnectedAt: rawEmailConnection.ConnectedAt,
	}, nil
}

// Get a list of the users connections
func (p *ProfileService) GetEmailConnections(userId string) ([]*EmailConnection, error) {
	rawEmailConnections, err := p.model.GetEmailConnections(userId)
	if err != nil {
		return nil, err
	}

	emailConnections := []*EmailConnection{}
	for _, rawEmailConnection := range rawEmailConnections {
		emailConnections = append(emailConnections, &EmailConnection{
			Id:          rawEmailConnection.Id,
			UserId:      rawEmailConnection.UserId,
			Email:       rawEmailConnection.Email,
			ConnectedAt: rawEmailConnection.ConnectedAt,
		})
	}

	return emailConnections, nil
}
