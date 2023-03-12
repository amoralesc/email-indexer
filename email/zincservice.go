package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	uploadPath      = "/api/_bulkv2"
	deleteIndexPath = "/api/index/emails"
)

func UploadEmails(bulk BulkEmails, zincUrl string, zincAdminUser string, zincAdminPassword string) error {
	// Convert the struct to JSON
	jsonBytes, err := json.Marshal(bulk)
	if err != nil {
		return err
	}
	jsonReader := bytes.NewReader(jsonBytes)

	// Create the post request
	req, err := http.NewRequest("POST", zincUrl+uploadPath, jsonReader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(zincAdminUser, zincAdminPassword)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}

func DeleteIndex(zincUrl string, zincAdminUser string, zincAdminPassword string) error {
	// Create the delete request
	req, err := http.NewRequest("DELETE", zincUrl+deleteIndexPath, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(zincAdminUser, zincAdminPassword)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}
