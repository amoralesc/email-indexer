package email

import (
	"bytes"
	"net/mail"
	"os"
)

type Email struct {
	MessageID string `json:"message-id"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Cc        string `json:"cc"`
	Bcc       string `json:"bcc"`
	Subject   string `json:"subject"`
	XFrom     string `json:"x-from"`
	XTo       string `json:"x-to"`
	XCc       string `json:"x-cc"`
	XBcc      string `json:"x-bcc"`
	Body      string `json:"body"`
}

// EmailFromFile converts an email file located at path
// to an Email struct for easy JSON encoding
func EmailFromFile(path string) (Email, error) {
	// Open the email file and read its contents
	file, err := os.Open(path)
	if err != nil {
		return Email{}, err
	}
	defer file.Close()

	// Parse the email
	msg, err := mail.ReadMessage(file)
	if err != nil {
		return Email{}, err
	}

	// Convert the msg to a struct
	emailStruct := Email{
		MessageID: msg.Header.Get("Message-ID"),
		Date:      msg.Header.Get("Date"),
		From:      msg.Header.Get("From"),
		To:        msg.Header.Get("To"),
		Cc:        msg.Header.Get("Cc"),
		Bcc:       msg.Header.Get("Bcc"),
		Subject:   msg.Header.Get("Subject"),
		XFrom:     msg.Header.Get("X-From"),
		XTo:       msg.Header.Get("X-To"),
		XCc:       msg.Header.Get("X-Cc"),
		XBcc:      msg.Header.Get("X-Bcc"),
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	emailStruct.Body = buf.String()

	return emailStruct, nil
}
