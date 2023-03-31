package zinc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const indexPath = "/api/index/"

// CheckIndex checks if the index exists in the zinc server
func (service *ZincService) CheckIndex() (bool, error) {
	// create the head request
	req, err := http.NewRequest("HEAD", service.Url+indexPath+"emails", nil)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(service.User, service.Password)

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return true, nil
}

// CreateIndex creates an index in the zinc server with a mapping that matches the Email struct
func (service *ZincService) CreateIndex() error {
	const emailsIndexMapping = `
	{
		"name": "emails",
		"storage_type": "disk",
		"mappings": {
			"properties": {
				"messageId": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
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
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"to": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
					"aggregatable": true,
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
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"bcc": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"body": {
					"type": "text",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"isRead": {
					"type": "boolean",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"isStarred": {
					"type": "boolean",
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
	req, err := http.NewRequest("POST", service.Url+indexPath, bytes.NewReader([]byte(emailsIndexMapping)))
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

// DeleteIndex deletes the emails index from the zinc server
func (service *ZincService) DeleteIndex() error {
	// create the delete request
	req, err := http.NewRequest("DELETE", service.Url+indexPath+"emails", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(service.User, service.Password)

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
