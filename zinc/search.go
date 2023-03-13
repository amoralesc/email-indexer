package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amoralesc/indexer/email"
)

const (
	defaultQuerySize = 100
	searchPath       = "/es/emails/_search"
)

type QuerySettings struct {
	From    int  // the offset to start from (pagination). Default: 0
	Size    int  // the number of elements to return (pagination). Default: 100
	SortAsc bool // if true, the elements will be sorted ascending by date. Default: false
}

// parseEmails parses the response from the zinc server
// and returns the emails
func parseEmails(body []byte) ([]email.Email, error) {
	// parse the response
	var resp struct {
		Hits struct {
			Hits []struct {
				Source email.Email `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// extract the emails
	return func() []email.Email {
		emails := make([]email.Email, len(resp.Hits.Hits))
		for i, hit := range resp.Hits.Hits {
			emails[i] = hit.Source
		}
		return emails
	}(), nil
}

// GetAllEmails returns all emails from the zinc server
func GetAllEmails(settings QuerySettings, serverAuth ServerAuth) ([]email.Email, error) {
	// create the query
	const query = `
	{
		"query": {
			"bool": {
				"must": [
					{ "match_all": {} }
				]
			}
		},
		"sort": [
			"%vdate"
		],
		"from": %d,
		"size": %d
	}
	`
	if settings.Size == 0 {
		settings.Size = defaultQuerySize
	}
	sort := "-"
	if settings.SortAsc {
		sort = "+"
	}
	queryStr := fmt.Sprintf(query, sort, settings.From, settings.Size)

	// create the request
	req, err := http.NewRequest("POST", serverAuth.Url+searchPath, bytes.NewBuffer([]byte(queryStr)))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(serverAuth.User, serverAuth.Password)

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

	// parse the response
	emails, err := parseEmails(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return emails, nil
}
