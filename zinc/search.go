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
	esSearchPath  = "/es/emails/_search"
	apiSearchPath = "/api/emails/_search"
)

// QueryResponse is the response from the zinc server to a query.
type QueryResponse struct {
	Total  int           `json:"total"`  // Total number of emails that match the query (not the number of emails returned)
	Took   int           `json:"took"`   // Time it took to execute the query
	Emails []email.Email `json:"emails"` // Emails that match the query (paginated)
}

// parseQueryResponse parses the body response from the zinc server
// into a QueryResponse struct.
func (service *ZincService) parseQueryResponse(body []byte) (*QueryResponse, error) {
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
		return nil, err
	}

	// extract the emails
	emails := func() []email.Email {
		emails := make([]email.Email, len(resp.Hits.Hits))
		for i, hit := range resp.Hits.Hits {
			emails[i] = hit.Source
		}
		return emails
	}()

	return &QueryResponse{
		Total:  resp.Hits.Total.Value,
		Took:   resp.Took,
		Emails: emails,
	}, nil
}

// sendQuery sends a query to the zinc server. It returns the emails that match the query.
func (service *ZincService) sendQuery(query string, searchPath string) (*QueryResponse, error) {
	// create the request
	req, err := http.NewRequest("POST", service.Url+searchPath, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(service.User, service.Password)

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
	queryResponse, err := service.parseQueryResponse(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return queryResponse, nil
}

// GetAllEmailAddresses returns all email addresses from the zinc server.
// This is a resource-intensive operation, so it should be used with caution.
func (service *ZincService) GetAllEmailAddresses() ([]string, error) {
	// create the query template
	const queryTemplate = `
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
		queryStr := fmt.Sprintf(queryTemplate, field, size)

		// create the request
		req, err := http.NewRequest("POST", service.Url+apiSearchPath, bytes.NewBuffer([]byte(queryStr)))
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(service.User, service.Password)

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
func (service *ZincService) GetAllEmails(settings *QuerySettings) (*QueryResponse, error) {
	// create the query template
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
	query := fmt.Sprintf(queryTemplate, settings.ParseQuerySettings())

	return service.sendQuery(query, esSearchPath)
}

// GetEmailByMessageId returns the email that has the given message id.
func (service *ZincService) GetEmailByMessageId(messageId string) (*QueryResponse, error) {
	// create the query template
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

	return service.sendQuery(query, esSearchPath)
}

// GetEmailsBySearchQuery returns all emails that match the given search query (paginated).
func (service *ZincService) GetEmailsBySearchQuery(searchQuery *SearchQuery, settings *QuerySettings) (*QueryResponse, error) {
	// create the query template
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

	query := fmt.Sprintf(queryTemplate, strings.Join(mustParameters, ", "), mustNotParameters, filterParameters, settings.ParseQuerySettings())

	return service.sendQuery(query, esSearchPath)
}

// GetEmailsByQueryString returns all emails that match the given query string (paginated).
// A query string is a string composed of query language syntax. For example:
// "query string +other word +content:test"
func (service *ZincService) GetEmailsByQueryString(queryString string, settings *QuerySettings) (*QueryResponse, error) {
	// create the query template
	const queryTemplate = `
	{
		"query": {
			"bool": {
				"must": [
					{ "query_string": { "query": "%v" } }
				]
			}
		},
		%v
	}
	`

	query := fmt.Sprintf(queryTemplate, queryString, settings.ParseQuerySettings())

	return service.sendQuery(query, esSearchPath)
}
