package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/amoralesc/indexer/email"
)

const (
	searchPath = "/es/emails/_search"
	aggPath    = "/api/emails/_search"
)

type QueryResponse struct {
	Total  int           `json:"total"`
	Took   int           `json:"took"`
	Emails []email.Email `json:"emails"`
}

// parseQueryResponse parses the query response from the zinc server.
func parseQueryResponse(body []byte) (QueryResponse, error) {
	// parse the response
	var resp struct {
		Took int `json:"took"`
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source email.Email `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return QueryResponse{}, err
	}

	// extract the emails
	emails := func() []email.Email {
		emails := make([]email.Email, len(resp.Hits.Hits))
		for i, hit := range resp.Hits.Hits {
			emails[i] = hit.Source
		}
		return emails
	}()

	return QueryResponse{
		Total:  resp.Hits.Total.Value,
		Took:   resp.Took,
		Emails: emails,
	}, nil
}

// sendQuery sends a query to the zinc server. It returns the emails that match the query.
func sendQuery(query string, serverAuth ServerAuth) (QueryResponse, error) {
	// create the request
	req, err := http.NewRequest("POST", serverAuth.Url+searchPath, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return QueryResponse{}, err
	}
	req.SetBasicAuth(serverAuth.User, serverAuth.Password)

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return QueryResponse{}, err
	}
	defer resp.Body.Close()

	// check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return QueryResponse{}, fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	// parse the response
	queryResponse, err := parseQueryResponse(body)
	if err != nil {
		return QueryResponse{}, fmt.Errorf("error parsing response: %v", err)
	}

	return queryResponse, nil
}

// GetAllEmailAddresses returns all email addresses from the zinc server.
// This is a resource-intensive operation, so it should be used with caution.
func GetAllEmailAddresses(serverAuth ServerAuth) ([]string, error) {
	// create the query
	const query = `
	{
		"search_type": "matchall",
		"max_results": 0,
		"aggs": {
			"results": {
				"agg_type": "term",
				"field": "%v",
				"size": %d
			}
		}
	}`
	const size = 100000 // this is a magic number, but it should be enough

	// the request is done for each field: from, to, cc, bcc
	var addresses map[string]struct{} // using a map to avoid duplicates
	for _, field := range []string{"from", "to", "cc", "bcc"} {
		queryStr := fmt.Sprintf(query, field, size)

		// create the request
		req, err := http.NewRequest("POST", serverAuth.Url+aggPath, bytes.NewBuffer([]byte(queryStr)))
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(serverAuth.User, serverAuth.Password)

		// send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// check the response
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
		}

		// parse the response
		var respStruct struct {
			Aggregations struct {
				Results struct {
					Buckets []struct {
						Key interface{} `json:"key"`
					} `json:"buckets"`
				} `json:"results"`
			} `json:"aggregations"`
		}
		if err := json.Unmarshal(body, &respStruct); err != nil {
			return nil, err
		}

		// extract the addresses
		for _, bucket := range respStruct.Aggregations.Results.Buckets {
			if addresses == nil {
				addresses = make(map[string]struct{})
			}
			// key can be a string or a number, if its a number it should be ignored
			if addr, ok := bucket.Key.(string); ok {
				addresses[addr] = struct{}{}
			}
		}
	}

	// convert the map to a slice
	addrs := make([]string, len(addresses))
	i := 0
	for addr := range addresses {
		addrs[i] = addr
		i++
	}

	return addrs, nil
}

// GetAllEmails returns all emails from the zinc server (paginated).
func GetAllEmails(settings QuerySettings, serverAuth ServerAuth) (QueryResponse, error) {
	// create the query
	const queryTemplate = `
	{
		"query": {
			"bool": {
				"must": [
					{ "match_all": {} }
				]
			}
		},
		%v
	}
	`
	if settings.Pagination.Size == 0 {
		settings.Pagination.Size = defaultQuerySize
	}
	query := fmt.Sprintf(queryTemplate, parseQuerySettings(settings))

	return sendQuery(query, serverAuth)
}

// GetEmailByMessageId returns the email that has the given message id.
func GetEmailByMessageId(messageId string, serverAuth ServerAuth) (QueryResponse, error) {
	// create the query
	const queryTemplate = `
	{
		"query": {
			"bool": {
				"must": [ %v ]
			}
		}
	}
	`
	query := fmt.Sprintf(queryTemplate, parseExactMatchParameter("message_id", messageId))

	resp, err := sendQuery(query, serverAuth)
	if err != nil {
		return QueryResponse{}, err
	}

	if len(resp.Emails) == 0 {
		return QueryResponse{}, fmt.Errorf("no email found with message id %v", messageId)
	}

	return resp, nil
}

// GetEmailsBySearchQuery returns all emails that match the given search query (paginated).
func GetEmailsBySearchQuery(searchQuery SearchQuery, settings QuerySettings, serverAuth ServerAuth) (QueryResponse, error) {
	// create the query
	const queryTemplate = `
	{
		"query": {
			"bool": {
				"must": [ %v ],
				"must_not": [ %v ],
				"filter": [ %v ]
			}
		},
		%v
	}
	`
	if settings.Pagination.Size == 0 {
		settings.Pagination.Size = defaultQuerySize
	}
	// parse the must parameters
	var mustParameters []string
	if searchQuery.From != "" {
		mustParameters = append(mustParameters, parseExactMatchParameter("from", searchQuery.From))
	}
	if len(searchQuery.To) > 0 {
		mustParameters = append(mustParameters, parseMultipleExactMatchParameter("to", searchQuery.To))
	}
	if len(searchQuery.Cc) > 0 {
		mustParameters = append(mustParameters, parseMultipleExactMatchParameter("cc", searchQuery.Cc))
	}
	if len(searchQuery.Bcc) > 0 {
		mustParameters = append(mustParameters, parseMultipleExactMatchParameter("bcc", searchQuery.Bcc))
	}
	if searchQuery.SubjectIncludes != "" {
		mustParameters = append(mustParameters, parseMatchTextParameter("subject", searchQuery.SubjectIncludes))
	}
	if searchQuery.BodyIncludes != "" {
		mustParameters = append(mustParameters, parseMatchTextParameter("body", searchQuery.BodyIncludes))
	}
	// parse the must_not parameters
	mustNotParameters := ""
	if searchQuery.BodyExcludes != "" {
		mustNotParameters = parseMatchTextParameter("body", searchQuery.BodyExcludes)
	}
	// parse the filter parameters
	filterParameters := parseDateRangeParameter(searchQuery.Date)

	// create the query string
	query := fmt.Sprintf(queryTemplate, strings.Join(mustParameters, ", "), mustNotParameters, filterParameters, parseQuerySettings(settings))

	return sendQuery(query, serverAuth)
}
