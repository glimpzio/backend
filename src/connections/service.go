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
	siteBaseUrl    string
}

// Create a new profile service
func NewConnectionService(readDb *sql.DB, writeDb *sql.DB, mailList *misc.MailList, profileService *profile.ProfileService, siteBaseUrl string) *ConnectionService {
	return &ConnectionService{model: &model.Model{ReadDb: readDb, WriteDb: writeDb}, mailList: mailList, profileService: profileService, siteBaseUrl: siteBaseUrl}
}

// Connect the users by email signup
func (c *ConnectionService) ConnectByEmail(inviteId string, email string, subscribe bool) (*CustomConnection, error) {
	_, user, err := c.profileService.GetInvite(inviteId)
	if err != nil {
		return nil, err
	}

	rawConnection, err := c.model.CreateCustomConnection(user.Id, nil, nil, nil, &email, nil, nil, nil)
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
	body += fmt.Sprintf("\n\nBy the way, if you want to boost your own leads and sales from networking events like %s, Glimpz makes it easy for you to connect and follow up with professionals at in-person networking events, helping you build long-lasting relationships. Try it for free now! %s?referral=uide-%s", user.FirstName, c.siteBaseUrl, user.Id)

	body += "\n\nWarm regards,\nBen"

	subject := fmt.Sprintf("Here's %s %s's Profile As You Requested!", user.FirstName, user.LastName)

	err = c.mailList.SendMail(email, email, subject, body)
	if err != nil {
		return nil, err
	}

	return &CustomConnection{
		Id:          rawConnection.Id,
		UserId:      rawConnection.UserId,
		ConnectedAt: rawConnection.ConnectedAt,
		FirstName:   rawConnection.FirstName,
		LastName:    rawConnection.LastName,
		Notes:       rawConnection.Notes,
		Email:       rawConnection.Email,
		Phone:       rawConnection.Phone,
		Website:     rawConnection.Website,
		LinkedIn:    rawConnection.LinkedIn,
	}, nil
}

// Upsert a custom connection
func (c *ConnectionService) UpsertCustomConnection(userId string, id *string, customConnection *NewCustomConnection) (*CustomConnection, error) {
	var rawConnection *model.CustomConnection
	var err error

	if id == nil {
		rawConnection, err = c.model.CreateCustomConnection(userId, customConnection.FirstName, customConnection.LastName, customConnection.Notes, customConnection.Email, customConnection.Phone, customConnection.Website, customConnection.LinkedIn)
	} else {
		existing, err := c.model.GetCustomConnection(*id)

		if err != nil {
			return nil, err
		} else if existing.UserId != userId {
			return nil, ErrNotAuthorized
		}

		rawConnection, err = c.model.UpdateCustomConnection(*id, customConnection.FirstName, customConnection.LastName, customConnection.Notes, customConnection.Email, customConnection.Phone, customConnection.Website, customConnection.LinkedIn)
	}

	if err != nil {
		return nil, err
	}

	return &CustomConnection{
		Id:          rawConnection.Id,
		UserId:      rawConnection.UserId,
		ConnectedAt: rawConnection.ConnectedAt,
		FirstName:   rawConnection.FirstName,
		LastName:    rawConnection.LastName,
		Notes:       rawConnection.Notes,
		Email:       rawConnection.Email,
		Phone:       rawConnection.Phone,
		Website:     rawConnection.Website,
		LinkedIn:    rawConnection.LinkedIn,
	}, nil
}

// Delete a custom connection
func (c *ConnectionService) DeleteCustomConnection(id string) (*CustomConnection, error) {
	rawConnection, err := c.model.DeleteCustomConnection(id)
	if err != nil {
		return nil, err
	}

	return &CustomConnection{
		Id:          rawConnection.Id,
		UserId:      rawConnection.UserId,
		ConnectedAt: rawConnection.ConnectedAt,
		FirstName:   rawConnection.FirstName,
		LastName:    rawConnection.LastName,
		Notes:       rawConnection.Notes,
		Email:       rawConnection.Email,
		Phone:       rawConnection.Phone,
		Website:     rawConnection.Website,
		LinkedIn:    rawConnection.LinkedIn,
	}, nil
}

// Get a custom connection
func (c *ConnectionService) GetCustomConnection(id string) (*CustomConnection, error) {
	rawConnection, err := c.model.GetCustomConnection(id)
	if err != nil {
		return nil, err
	}

	return &CustomConnection{
		Id:          rawConnection.Id,
		UserId:      rawConnection.UserId,
		ConnectedAt: rawConnection.ConnectedAt,
		FirstName:   rawConnection.FirstName,
		LastName:    rawConnection.LastName,
		Notes:       rawConnection.Notes,
		Email:       rawConnection.Email,
		Phone:       rawConnection.Phone,
		Website:     rawConnection.Website,
		LinkedIn:    rawConnection.LinkedIn,
	}, nil
}

// Get a list of the users custom connections
func (c *ConnectionService) GetCustomConnections(userId string) ([]*CustomConnection, error) {
	rawConnections, err := c.model.GetCustomConnections(userId)
	if err != nil {
		return nil, err
	}

	connections := []*CustomConnection{}
	for _, rawConnection := range rawConnections {
		connections = append(connections, &CustomConnection{
			Id:          rawConnection.Id,
			UserId:      rawConnection.UserId,
			ConnectedAt: rawConnection.ConnectedAt,
			FirstName:   rawConnection.FirstName,
			LastName:    rawConnection.LastName,
			Notes:       rawConnection.Notes,
			Email:       rawConnection.Email,
			Phone:       rawConnection.Phone,
			Website:     rawConnection.Website,
			LinkedIn:    rawConnection.LinkedIn,
		})
	}

	return connections, nil
}
