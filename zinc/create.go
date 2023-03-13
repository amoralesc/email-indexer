package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amoralesc/indexer/email"
)

// BulkEmails is used to upload emails in bulk
// to the zinc server
type BulkEmails struct {
	Index   string        `json:"index"`
	Records []email.Email `json:"records"`
}

const (
	uploadPath = "/api/_bulkv2"
	indexPath  = "/api/index/"
)

// CreateIndex creates an index in the zinc server
// with a mapping that matches the email.Email struct
func CreateIndex(serverAuth ServerAuth) error {
	const emailsIndexMapping = `
	{
		"name": "emails",
		"storage_type": "disk",
		"mappings": {
			"properties": {
				"message-id": {
					"type": "keyword",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"date": {
					"type": "date",
					"format": "2006-01-02T15:04:05Z07:00",
					"index": true,
					"store": false,
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"from": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"to": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"subject": {
					"type": "text",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"cc": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"bcc": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"body": {
					"type": "text",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				}
			}
		}
	}`

	// create the post request
	req, err := http.NewRequest("POST", serverAuth.Url+indexPath, bytes.NewReader([]byte(emailsIndexMapping)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(serverAuth.User, serverAuth.Password)
	req.Header.Set("Content-Type", "application/json")

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}

// UploadEmails uploads a list of emails to the zinc server
func UploadEmails(bulk BulkEmails, serverAuth ServerAuth) error {
	// convert the struct to JSON
	jsonBytes, err := json.Marshal(bulk)
	if err != nil {
		return err
	}

	// create the post request
	req, err := http.NewRequest("POST", serverAuth.Url+uploadPath, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.SetBasicAuth(serverAuth.User, serverAuth.Password)
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

// DeleteIndex deletes the emails index from the zinc server
func DeleteIndex(serverAuth ServerAuth) error {
	// create the delete request
	req, err := http.NewRequest("DELETE", serverAuth.Url+indexPath+"emails", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(serverAuth.User, serverAuth.Password)

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}
