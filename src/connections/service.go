package connections

import (
	"database/sql"
	"fmt"

	"github.com/glimpzio/backend/connections/model"
	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile"
)

type ConnectionService struct {
	model          *model.Model
	mailList       *misc.MailList
	profileService *profile.ProfileService
}

// Create a new profile service
func NewConnectionService(db *sql.DB, mailList *misc.MailList, profileService *profile.ProfileService) *ConnectionService {
	return &ConnectionService{model: &model.Model{Db: db}, mailList: mailList, profileService: profileService}
}

// Connec the users by email signup
func (c *ConnectionService) ConnectByEmail(inviteId string, email string, subscribe bool) (*EmailConnection, error) {
	_, user, err := c.profileService.GetInvite(inviteId)
	if err != nil {
		return nil, err
	}

	rawEmailConnection, err := c.model.CreateEmailConnection(user.Id, email)
	if err != nil {
		return nil, err
	}

	if subscribe {
		err = c.mailList.AddMarketing(email, nil, nil)
		if err != nil {
			return nil, err
		}
	}

	body := fmt.Sprintf("Hey, hope you're well!\n\nAs you requested, here's the Glimpz profile for %s %s:\n\n- Bio: %s", user.FirstName, user.LastName, user.Bio)

	if user.Profile.Email != nil {
		body += fmt.Sprintf("\n- Email: %s", *user.Profile.Email)
	}
	if user.Profile.Phone != nil {
		body += fmt.Sprintf("\n- Phone: %s", *user.Profile.Phone)
	}
	if user.Profile.Website != nil {
		body += fmt.Sprintf("\n- Website: %s", *user.Profile.Website)
	}
	if user.Profile.LinkedIn != nil {
		body += fmt.Sprintf("\n- LinkedIn: %s", *user.Profile.LinkedIn)
	}

	body += fmt.Sprintf("\n\nWe've also forwarded your email to %s so they can follow up with you when they get a chance.", user.FirstName)
	body += fmt.Sprintf("\n\nBy the way, if you're ever looking to increase your own leads and sales from networking events like %s, did you know that you can make your own Glimpz profile for free right now? Glimpz makes it easy for you to connect with other professionals at networking events and convert them into long-lasting business partners or clients. Check it out at https://glimpz.io?referral=uide-%s", user.FirstName, user.Id)

	body += "\n\nWarm regards,\nBen"

	subject := fmt.Sprintf("Here's %s %s's Profile As You Requested!", user.FirstName, user.LastName)

	err = c.mailList.SendMail(email, email, subject, body)
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
func (c *ConnectionService) GetEmailConnections(userId string) ([]*EmailConnection, error) {
	rawEmailConnections, err := c.model.GetEmailConnections(userId)
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
