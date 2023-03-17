package zinc

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	defaultQueryStart = 0
	defaultQuerySize  = 100
	defaultSortField  = "-date"
)

// QueryPaginationSettings sets the pagination parameters for the query.
type QueryPaginationSettings struct {
	Start int // the offset to start from (pagination). Default: 0
	Size  int // the number of elements to return (pagination). Default: 100
}

// QuerySettings sets the parameters for the query.
type QuerySettings struct {
	Sort       string                   // the sorting parameters. Default: "-date"
	Pagination *QueryPaginationSettings // the pagination parameters. Default: {Start: 0, Size: 100}
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

// ValidateSortField validates a sort field with the format: (+|-)(from|to|cc|bcc|date)
func ValidateSortField(sortField string) error {
	matches := regexp.MustCompile(`^(\+|-)(from|to|cc|bcc|date)$`).MatchString(sortField)
	if !matches {
		return fmt.Errorf("invalid sort field: %v", sortField)
	}
	return nil
}

// NewQuerySettings creates new QuerySettings.
func NewQuerySettings(sortBy string, page, pageSize int) (*QuerySettings, error) {
	sortFields := strings.Split(sortBy, ",")
	for _, s := range sortFields {
		err := ValidateSortField(s)
		if err != nil {
			return nil, err
		}
	}

	if sortBy == "" {
		sortBy = defaultSortField
	}
	if page < 0 {
		return nil, fmt.Errorf("page should be equal or greater than 0: %v", page)
	}
	if pageSize < 0 {
		return nil, fmt.Errorf("pageSize should be equal or greater than 0: %v", pageSize)
	}
	if pageSize == 0 {
		pageSize = defaultQuerySize
	}

	return &QuerySettings{Sort: sortBy, Pagination: &QueryPaginationSettings{Start: page, Size: pageSize}}, nil
}

// parseQuerySettings parses the query settings to a string.
func (settings *QuerySettings) ParseQuerySettings() string {
	return fmt.Sprintf(`"sort": [ %v ], "from": %d, "size": %d`,
		settings.Sort,
		settings.Pagination.Start,
		settings.Pagination.Size)
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
