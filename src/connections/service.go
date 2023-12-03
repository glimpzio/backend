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
	landingBaseUrl string
}

// Create a new profile service
func NewConnectionService(readDb *sql.DB, writeDb *sql.DB, mailList *misc.MailList, profileService *profile.ProfileService, landingBaseUrl string) *ConnectionService {
	return &ConnectionService{model: &model.Model{ReadDb: readDb, WriteDb: writeDb}, mailList: mailList, profileService: profileService, landingBaseUrl: landingBaseUrl}
}

// Connect the users by email signup
func (c *ConnectionService) Connect(inviteId string, subscribe bool, email string, firstName *string, lastName *string) (*CustomConnection, error) {
	_, user, err := c.profileService.GetInvite(inviteId)
	if err != nil {
		return nil, err
	}

	rawConnection, err := c.model.CreateCustomConnection(user.Id, firstName, lastName, nil, &email, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if subscribe {
		err = c.mailList.AddMarketing(email, firstName, lastName)
		if err != nil {
			return nil, err
		}
	}

	body := ""

	if firstName != nil {
		body += fmt.Sprintf("Hey %s, hope you're well!", *firstName)
	} else {
		body += "Hey, hope you're well!"
	}

	body += fmt.Sprintf("\n\nAs you requested, here's the Glimpz profile for %s %s:\n\n- Bio: %s", user.FirstName, user.LastName, user.Bio)

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
	body += fmt.Sprintf("\n\nBy the way, if you want to boost your own leads and sales from networking events like %s, Glimpz makes it easy for you to connect and follow up with professionals at in-person networking events, helping you build long-lasting relationships. Try it for free now! %s?referral=uide-%s", user.FirstName, c.landingBaseUrl, user.Id)

	body += "\n\nWarm regards,\nBen"

	subject := fmt.Sprintf("Here's %s %s's Profile As You Requested!", user.FirstName, user.LastName)

	var name string

	if firstName != nil {
		name = *firstName
	} else {
		name = email
	}

	err = c.mailList.SendMail(name, email, subject, body)
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
func (c *ConnectionService) UpsertCustomConnection(userId string, id *string, newCustomConnection *NewCustomConnection) (*CustomConnection, error) {
	var rawConnection *model.CustomConnection

	if id == nil {
		connection, err := c.model.CreateCustomConnection(userId, newCustomConnection.FirstName, newCustomConnection.LastName, newCustomConnection.Notes, newCustomConnection.Email, newCustomConnection.Phone, newCustomConnection.Website, newCustomConnection.LinkedIn)
		if err != nil {
			return nil, err
		}

		rawConnection = connection
	} else {
		existing, err := c.model.GetCustomConnection(*id)

		if err != nil {
			return nil, err
		} else if existing.UserId != userId {
			return nil, ErrNotAuthorized
		}

		connection, err := c.model.UpdateCustomConnection(*id, newCustomConnection.FirstName, newCustomConnection.LastName, newCustomConnection.Notes, newCustomConnection.Email, newCustomConnection.Phone, newCustomConnection.Website, newCustomConnection.LinkedIn)
		if err != nil {
			return nil, err
		}

		rawConnection = connection
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
	} else if rawConnection == nil {
		return nil, ErrDoesNotExist
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
func (c *ConnectionService) GetCustomConnections(userId string, limit int, offset int) ([]*CustomConnection, error) {
	rawConnections, err := c.model.GetCustomConnections(userId, limit, offset)
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
