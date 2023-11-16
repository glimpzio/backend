package misc

import (
	"encoding/json"
	"errors"

	"github.com/sendgrid/sendgrid-go"
)

type MailList struct {
	sendgridKey string
	listId      string
}

// Create a mail client
func NewMailList(sendgridKey string, listId string) *MailList {
	return &MailList{sendgridKey: sendgridKey, listId: listId}
}

// Add an email to a mailing list
func (m *MailList) Add(name string, email string) error {

	contactData := map[string]interface{}{
		"list_ids": []string{m.listId},
		"contacts": []interface{}{
			map[string]interface{}{
				"email":         email,
				"custom_fields": map[string]string{"name": name},
			},
		},
	}

	body, err := json.Marshal(contactData)
	if err != nil {
		return err
	}

	request := sendgrid.GetRequest(m.sendgridKey, "/v3/marketing/contacts", "https://api.sendgrid.com")
	request.Method = "PUT"
	request.Headers["Content-Type"] = "application/json"
	request.Body = body

	response, err := sendgrid.API(request)
	if err != nil {
		return err
	} else if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New("invalid email response code")
	}

	return nil
}
