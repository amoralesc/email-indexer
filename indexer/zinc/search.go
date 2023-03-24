package zinc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	esSearchPath    = "/es/emails/_search"
	apiDocumentPath = "/api/emails/_doc"
)

// EmailWithId is the returned email format from the zinc server.
type EmailWithId struct {
	Id        string    `json:"_id"`
	MessageId string    `json:"messageId"`
	Date      time.Time `json:"date"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Cc        []string  `json:"cc"`
	Bcc       []string  `json:"bcc"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	IsRead    bool      `json:"isRead"`
	IsStarred bool      `json:"isStarred"`
}

// QueryResponse is the response from the zinc server to a query.
type QueryResponse struct {
	Total  int           `json:"total"`  // Total number of emails that match the query (not the number of emails returned)
	Took   int           `json:"took"`   // Time it took to execute the query
	Emails []EmailWithId `json:"emails"` // Emails that match the query (paginated)
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
				Id     string      `json:"_id"`
				Source EmailWithId `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// extract the emails
	emails := func() []EmailWithId {
		emails := make([]EmailWithId, len(resp.Hits.Hits))
		for i, hit := range resp.Hits.Hits {
			emails[i] = hit.Source
			emails[i].Id = hit.Id
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
func (service *ZincService) sendQuery(query string) (*QueryResponse, error) {
	// create the request
	req, err := http.NewRequest("POST", service.Url+esSearchPath, bytes.NewBuffer([]byte(query)))
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
				%v
			}
		},
		%v
	}
	`

	var filter = `, "filter": [ %v ]`
	if settings.StarredOnly {
		filter = fmt.Sprintf(filter, settings.ParseStarredFilter())
	} else {
		filter = ""
	}

	query := fmt.Sprintf(queryTemplate, filter, settings.ParseQuerySettings())

	return service.sendQuery(query)
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
	var filterParameters []string
	filterParameters = append(filterParameters, parseDateRangeParameter(searchQuery.DateRange))
	if settings.StarredOnly {
		filterParameters = append(filterParameters, settings.ParseStarredFilter())
	}

	query := fmt.Sprintf(queryTemplate, strings.Join(mustParameters, ", "), mustNotParameters, strings.Join(filterParameters, ", "), settings.ParseQuerySettings())

	return service.sendQuery(query)
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
				%v
			}
		},
		%v
	}
	`

	var filter = `, "filter": [ %v ]`
	if settings.StarredOnly {
		filter = fmt.Sprintf(filter, settings.ParseStarredFilter())
	} else {
		filter = ""
	}

	query := fmt.Sprintf(queryTemplate, queryString, filter, settings.ParseQuerySettings())

	return service.sendQuery(query)
}

// GetEmailByMessageId returns the email that has the given message id.
func (service *ZincService) GetEmailByMessageId(messageId string) (*EmailWithId, error) {
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

	queryResponse, err := service.sendQuery(query)
	if err != nil {
		return nil, err
	}
	if len(queryResponse.Emails) == 0 {
		return nil, fmt.Errorf("messsage id not found %v", messageId)
	}

	return &queryResponse.Emails[0], nil
}

// GetEmailById returns the email that has the given _id (zinc id).
func (service *ZincService) GetEmailById(id string) (*EmailWithId, error) {
	// create the request
	req, err := http.NewRequest("GET", service.Url+apiDocumentPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(service.User, service.Password)

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
	var respStruct struct {
		Id     string      `json:"_id"`
		Source EmailWithId `json:"_source"`
	}

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	email := respStruct.Source
	email.Id = respStruct.Id

	return &email, nil
}
