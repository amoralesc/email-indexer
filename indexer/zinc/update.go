package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amoralesc/email-indexer/indexer/email"
)

const apiUpdatePath = "/api/emails/_update"

// UpdateEmail updates an email in the zinc server.
func (service *ZincService) UpdateEmail(id string, email *email.Email) (*EmailResponse, error) {
	jsonBytes, err := json.Marshal(*email)
	if err != nil {
		return nil, err
	}

	// create the request
	req, err := http.NewRequest("POST", service.Url+apiUpdatePath+"/"+id, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(service.User, service.Password)

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	// Create emailResponse from email
	emailResponse := EmailResponse{
		Id:        id,
		MessageId: email.MessageId,
		Date:      email.Date,
		From:      email.From,
		To:        email.To,
		Cc:        email.Cc,
		Bcc:       email.Bcc,
		Subject:   email.Subject,
		Body:      email.Body,
		IsRead:    email.IsRead,
		IsStarred: email.IsStarred,
	}

	return &emailResponse, nil
}
