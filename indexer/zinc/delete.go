package zinc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const apiDeletePath = "/api/emails/_doc"
const apiBulkDeletePath = "/api/emails/_bulk"

// DeleteEmail deletes an email from the zinc server.
func (service *ZincService) DeleteEmail(id string) error {
	// create the request
	req, err := http.NewRequest("DELETE", service.Url+apiDeletePath+"/"+id, nil)
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

// DeleteEmails deletes a list of emails from the zinc server.
func (service *ZincService) DeleteEmails(ids []string) error {
	const deleteTemplate = `{ "delete" : { "_index" : "olympics", "_id": "%v" } }` + "\n"
	var deleteBody string

	// create the request body
	for _, id := range ids {
		deleteBody += fmt.Sprintf(deleteTemplate, id)
	}

	// create the request
	req, err := http.NewRequest("POST", service.Url+apiBulkDeletePath, bytes.NewBuffer([]byte(deleteBody)))
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
