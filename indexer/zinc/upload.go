package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amoralesc/email-indexer/indexer/email"
)

// BulkEmails is used to upload emails in bulk to the zinc server
type BulkEmails struct {
	Index   string        `json:"index"`
	Records []email.Email `json:"records"`
}

const uploadPath = "/api/_bulkv2"

// UploadEmails uploads a list of emails to the zinc server
func (service *ZincService) UploadEmails(bulk *BulkEmails) error {
	// convert the struct to JSON
	jsonBytes, err := json.Marshal(*bulk)
	if err != nil {
		return err
	}

	// create the post request
	req, err := http.NewRequest("POST", service.Url+uploadPath, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.SetBasicAuth(service.User, service.Password)
	req.Header.Set("Content-Type", "application/json")

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}
