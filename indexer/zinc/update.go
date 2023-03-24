package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/amoralesc/email-indexer/indexer/email"
)

const apiUpdatePath = "/api/emails/_update"
const apiMultiUpdatePath = "/api/emails/_multi"

// UpdateEmail updates an email in the zinc server.
func (service *ZincService) UpdateEmail(id string, email *email.Email) (*EmailWithId, error) {
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

	// Create EmailWithId from email
	EmailWithId := EmailWithId{
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

	return &EmailWithId, nil
}

// UpdateEmails updates a list of emails in the zinc server.
func (service *ZincService) UpdateEmails(emails []*EmailWithId) ([]*EmailWithId, error) {
	jsonBytes, err := json.Marshal(emails)
	if err != nil {
		return nil, err
	}

	// remove [ ] (first and last char)
	jsonBytes = jsonBytes[1 : len(jsonBytes)-1]
	// replace all }, with }\n
	jsonBytes = bytes.ReplaceAll(jsonBytes, []byte("},"), []byte("}\n"))

	// print the json
	log.Println(string(jsonBytes))

	// create the request
	req, err := http.NewRequest("POST", service.Url+apiMultiUpdatePath, bytes.NewReader(jsonBytes))

	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(service.User, service.Password)

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return emails, nil
}
