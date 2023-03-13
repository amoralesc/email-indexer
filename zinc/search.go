package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/amoralesc/indexer/email"
)

const (
	defaultQuerySize = 100
	defaultSortField = "date"
	searchPath       = "/es/emails/_search"
	aggPath          = "/api/emails/_search"
)

type SortQuerySettings struct {
	Field   string // the field to sort by
	SortAsc bool   // if true, the elements will be sorted ascending. Default: false
}

// QuerySettings sets basic parameters for the query (pagination, sorting).
type QuerySettings struct {
	Start int                 // the offset to start from (pagination). Default: 0
	Size  int                 // the number of elements to return (pagination). Default: 100
	Sort  []SortQuerySettings // the fields to sort by in order of priority. Default: date descending
}

// DateRange represents a range of dates (from, to) to filter the query.
type DateRange struct {
	From time.Time `json:"from"` // the start date. Default: 0
	To   time.Time `json:"to"`   // the end date. Default: max time
}

// SearchQuery represents a query to search for emails.
// The query will only return emails that match all the fields.
// If a field is empty, it will be ignored.
type SearchQuery struct {
	From            string    `json:"from"`             // from address (exact match)
	To              []string  `json:"to"`               // to addresses (exact match to all)
	Cc              []string  `json:"cc"`               // cc addresses (exact match to all)
	Bcc             []string  `json:"bcc"`              // bcc addresses (exact match to all)
	SubjectIncludes string    `json:"subject_includes"` // subject (has text)
	BodyIncludes    string    `json:"body_includes"`    // body includes (has text)
	BodyExcludes    string    `json:"body_excludes"`    // body excludes (does not have text)
	Date            DateRange `json:"date"`             // the date range to filter the query
}

// parseSortQuerySettings parses the sort settings to a string.
func parseSortQuerySettings(sort []SortQuerySettings) string {
	if sort == nil {
		return `"-date"`
	}
	var sortStr string
	for _, s := range sort {
		if sortStr != `` {
			sortStr += `, `
		}
		if s.SortAsc {
			sortStr += `+`
		} else {
			sortStr += `-`
		}
		sortStr += `"` + s.Field + `"`
	}
	return sortStr
}

// parseEmails parses the query response from the zinc server
// and returns the emails.
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

func queryEmails(query string, serverAuth ServerAuth) ([]email.Email, error) {
	// create the request
	req, err := http.NewRequest("POST", serverAuth.Url+searchPath, bytes.NewBuffer([]byte(query)))
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
	emails, err := parseEmails(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return emails, nil
}

func parseSearchParameter(searchType string, field string, value string) string {
	return fmt.Sprintf(`{ "%v": { "%v": "%v" } }`, searchType, field, value)
}

func parseExactMatchParameter(field string, value string) string {
	return parseSearchParameter("term", field, value)
}

func parseMultipleExactMatchParameter(field string, values []string) string {
	parameters := make([]string, len(values))
	for i, value := range values {
		parameters[i] = parseExactMatchParameter(field, value)
	}
	return strings.Join(parameters, ", ")
}

func parseMatchTextParameter(field string, value string) string {
	return parseSearchParameter("match", field, value)
}

func parseDateRangeParameter(date DateRange) string {
	const rangeTemplate = `{ "range": { "date": { "format": "%v", %v } } }`

	fromTo := fmt.Sprintf(`"gte": "%v"`, date.From.Format(time.RFC3339))
	if !date.To.IsZero() {
		fromTo += fmt.Sprintf(`, "lte": "%v"`, date.To.Format(time.RFC3339))
	}

	return fmt.Sprintf(rangeTemplate, time.RFC3339, fromTo)
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
func GetAllEmails(settings QuerySettings, serverAuth ServerAuth) ([]email.Email, error) {
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
		"sort": [ %v ],
		"from": %d,
		"size": %d
	}
	`
	if settings.Size == 0 {
		settings.Size = defaultQuerySize
	}
	query := fmt.Sprintf(queryTemplate, parseSortQuerySettings(settings.Sort), settings.Start, settings.Size)

	return queryEmails(query, serverAuth)
}

// GetEmailByMessageId returns the email that has the given message id.
func GetEmailByMessageId(messageId string, serverAuth ServerAuth) (email.Email, error) {
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

	emails, err := queryEmails(query, serverAuth)
	if err != nil {
		return email.Email{}, err
	}

	if len(emails) == 0 {
		return email.Email{}, fmt.Errorf("email not found")
	}

	return emails[0], nil
}

// GetEmailsBySearchQuery returns all emails that match the given search query (paginated).
func GetEmailsBySearchQuery(searchQuery SearchQuery, settings QuerySettings, serverAuth ServerAuth) ([]email.Email, error) {
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
		"sort": [ %v ],
		"from": %d,
		"size": %d
	}
	`
	if settings.Size == 0 {
		settings.Size = defaultQuerySize
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
	query := fmt.Sprintf(queryTemplate, strings.Join(mustParameters, ", "), mustNotParameters, filterParameters, parseSortQuerySettings(settings.Sort), settings.Start, settings.Size)

	return queryEmails(query, serverAuth)
}
