package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amoralesc/indexer/email"
)

type BulkEmails struct {
	Index   string        `json:"index"`
	Records []email.Email `json:"records"`
}

const (
	uploadPath = "/api/_bulkv2"
	indexPath  = "/api/index/"
)

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

func CreateIndex(server ServerAuth) error {
	// create the post request
	req, err := http.NewRequest("POST", server.Url+indexPath, bytes.NewReader([]byte(emailsIndexMapping)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(server.User, server.Password)
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

func UploadEmails(bulk BulkEmails, server ServerAuth) error {
	// convert the struct to JSON
	jsonBytes, err := json.Marshal(bulk)
	if err != nil {
		return err
	}

	// create the post request
	req, err := http.NewRequest("POST", server.Url+uploadPath, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.SetBasicAuth(server.User, server.Password)
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

func DeleteIndex(indexName string, server ServerAuth) error {
	// create the delete request
	req, err := http.NewRequest("DELETE", server.Url+indexPath+indexName, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(server.User, server.Password)

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
