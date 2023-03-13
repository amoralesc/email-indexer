package email

import (
	"bytes"
	"net/mail"
	"os"
	"time"
)

// Email represents an email message that can
// be JSON encoded
type Email struct {
	MessageID string    `json:"message_id"`
	Date      time.Time `json:"date"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Cc        []string  `json:"cc"`
	Bcc       []string  `json:"bcc"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
}

// EmailFromFile parses an email file located at path
// to an Email struct for easy JSON encoding
func EmailFromFile(path string) (Email, error) {
	// open the email file and read its contents
	file, err := os.Open(path)
	if err != nil {
		return Email{}, err
	}
	defer file.Close()

	// parse the email
	msg, err := mail.ReadMessage(file)
	if err != nil {
		return Email{}, err
	}

	// convert the msg to a struct
	emailStruct := Email{
		MessageID: msg.Header.Get("Message-ID"),
		From:      msg.Header.Get("From"),
		Subject:   msg.Header.Get("Subject"),
	}
	// parse the date
	date, err := msg.Header.Date()
	if err != nil {
		return Email{}, err
	}
	emailStruct.Date = date
	// parse the To header if it exists
	if msg.Header.Get("To") != "" {
		to, err := msg.Header.AddressList("To")
		if err != nil {
			return Email{}, err
		}
		for _, addr := range to {
			emailStruct.To = append(emailStruct.To, addr.Address)
		}
	}
	// parse the Cc header if it exists
	if msg.Header.Get("Cc") != "" {
		cc, err := msg.Header.AddressList("Cc")
		if err != nil {
			return Email{}, err
		}
		for _, addr := range cc {
			emailStruct.Cc = append(emailStruct.Cc, addr.Address)
		}
	}
	// parse the Bcc header if it exists
	if msg.Header.Get("Bcc") != "" {
		bcc, err := msg.Header.AddressList("Bcc")
		if err != nil {
			return Email{}, err
		}
		for _, addr := range bcc {
			emailStruct.Bcc = append(emailStruct.Bcc, addr.Address)
		}
	}
	// parse the body
	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	emailStruct.Body = buf.String()

	return emailStruct, nil
}
