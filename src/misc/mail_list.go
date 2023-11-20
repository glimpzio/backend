package misc

import (
	"encoding/json"
	"errors"

	"github.com/sendgrid/sendgrid-go"
)

type MailList struct {
	sendgridKey     string
	accountListId   string
	marketingListId string
}

// Create a mail client
func NewMailList(sendgridKey string, accountlistId string, marketingListId string) *MailList {
	return &MailList{sendgridKey: sendgridKey, accountListId: accountlistId, marketingListId: marketingListId}
}

// Add to a mailing list
func (m *MailList) add(data *map[string]interface{}) error {
	body, err := json.Marshal(data)
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

// Add an email to the account list
func (m *MailList) AddAccount(email string, firstName string, lastName string) error {
	contactData := &map[string]interface{}{
		"list_ids": []string{m.accountListId},
		"contacts": []interface{}{
			map[string]interface{}{
				"email":      email,
				"first_name": firstName,
				"last_name":  lastName,
			},
		},
	}

	return m.add(contactData)
}

// Add an email to the marketing list
func (m *MailList) AddMarketing(email string, firstName *string, lastName *string) error {
	contactData := &map[string]interface{}{
		"list_ids": []string{m.marketingListId},
		"contacts": []interface{}{
			map[string]interface{}{
				"email":      email,
				"first_name": firstName,
				"last_name":  lastName,
			},
		},
	}

	return m.add(contactData)
}
