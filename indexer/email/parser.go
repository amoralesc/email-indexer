package email

import (
	"bytes"
	"net/mail"
	"os"
	"time"
)

// Email represents an email message that can be JSON encoded.
type Email struct {
	MessageID string    `json:"message_id"`
	Date      time.Time `json:"date"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Cc        []string  `json:"cc"`
	Bcc       []string  `json:"bcc"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	IsRead    bool      `json:"is_read"`
	IsStarred bool      `json:"is_starred"`
}

// EmailFromFile parses an email file located at path to an Email struct for easy JSON encoding.
func EmailFromFile(path string) (*Email, error) {
	// open the email file and read its contents
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// parse the email
	msg, err := mail.ReadMessage(file)
	if err != nil {
		return nil, err
	}

	// convert the msg to a struct
	emailObj := &Email{
		MessageID: msg.Header.Get("Message-ID"),
		From:      msg.Header.Get("From"),
		Subject:   msg.Header.Get("Subject"),
		IsRead:    false,
		IsStarred: false,
	}
	// parse the date
	date, err := msg.Header.Date()
	if err != nil {
		return nil, err
	}
	emailObj.Date = date
	// parse the To header if it exists
	if msg.Header.Get("To") != "" {
		to, err := msg.Header.AddressList("To")
		if err != nil {
			return nil, err
		}
		for _, addr := range to {
			emailObj.To = append(emailObj.To, addr.Address)
		}
	}
	// parse the Cc header if it exists
	if msg.Header.Get("Cc") != "" {
		cc, err := msg.Header.AddressList("Cc")
		if err != nil {
			return nil, err
		}
		for _, addr := range cc {
			emailObj.Cc = append(emailObj.Cc, addr.Address)
		}
	}
	// parse the Bcc header if it exists
	if msg.Header.Get("Bcc") != "" {
		bcc, err := msg.Header.AddressList("Bcc")
		if err != nil {
			return nil, err
		}
		for _, addr := range bcc {
			emailObj.Bcc = append(emailObj.Bcc, addr.Address)
		}
	}
	// parse the body
	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	emailObj.Body = buf.String()

	return emailObj, nil
}
