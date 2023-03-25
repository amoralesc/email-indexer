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
	defaultSortFields = "-date,messageId"
)

// QueryPaginationSettings sets the pagination parameters for the query.
type QueryPaginationSettings struct {
	Start int // the offset to start from (pagination). Default: 0
	Size  int // the number of elements to return (pagination). Default: 100
}

// QuerySettings sets the parameters for the query.
type QuerySettings struct {
	Sort        string                   // the sorting parameters. Default: "-date"
	Pagination  *QueryPaginationSettings // the pagination parameters. Default: {Start: 0, Size: 100}
	StarredOnly bool                     // if true, only starred emails will be returned. Default: false
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
	From            string    `json:"from"`            // from address (exact match)
	To              []string  `json:"to"`              // to addresses (exact match to all)
	Cc              []string  `json:"cc"`              // cc addresses (exact match to all)
	Bcc             []string  `json:"bcc"`             // bcc addresses (exact match to all)
	SubjectIncludes string    `json:"subjectIncludes"` // subject (has text)
	BodyIncludes    string    `json:"bodyIncludes"`    // body includes (has text)
	BodyExcludes    string    `json:"bodyExcludes"`    // body excludes (does not have text)
	DateRange       DateRange `json:"dateRange"`       // the date range to filter the query
}

// ValidateSortField validates a sort field with the format: (+|-)(from|to|cc|bcc|date)
func ValidateSortField(sortField string) error {
	matches := regexp.MustCompile(`^-?(messageId|date|from|to|cc|bcc)$`).MatchString(sortField)
	if !matches {
		return fmt.Errorf("invalid sort field: %v", sortField)
	}
	return nil
}

// NewQuerySettings creates new QuerySettings.
func NewQuerySettings(sortBy string, start, size int, starredOnly bool) (*QuerySettings, error) {
	if sortBy == "" {
		sortBy = defaultSortFields
	} else {
		// add default sort fields at end if not already present
		// this ensures the sort order is always the same
		// for the fields specified
		for _, s := range strings.Split(defaultSortFields, ",") {
			if !strings.Contains(sortBy, s) {
				sortBy += "," + s
			}
		}
	}

	sortFields := strings.Split(sortBy, ",")
	for _, s := range sortFields {
		err := ValidateSortField(s)
		if err != nil {
			return nil, err
		}
	}

	if start < 0 {
		return nil, fmt.Errorf("start should be equal or greater than 0: %v", start)
	}
	if size < 0 {
		return nil, fmt.Errorf("size should be equal or greater than 0: %v", size)
	}
	if size == 0 {
		size = defaultQuerySize
	}

	return &QuerySettings{Sort: sortBy, Pagination: &QueryPaginationSettings{Start: start, Size: size}, StarredOnly: starredOnly}, nil
}

// ParseQuerySortSettings parses the query sort settings to a string.
// (only pagination and sort since starred is a filter)3
func (settings *QuerySettings) ParseQuerySortSettings() string {
	sortFields := strings.Split(settings.Sort, ",")
	sortFieldsStr := make([]string, len(sortFields))
	for i, s := range sortFields {
		if !strings.HasPrefix(s, "-") {
			sortFieldsStr[i] = fmt.Sprintf(`"+%v"`, s)
		} else {
			sortFieldsStr[i] = fmt.Sprintf(`"%v"`, s)
		}
	}
	return strings.Join(sortFieldsStr, ",")
}

// ParseQuerySettings parses the query settings to a string.
func (settings *QuerySettings) ParseQuerySettings() string {
	return fmt.Sprintf(`"sort": [ %v ], "from": %d, "size": %d`,
		settings.ParseQuerySortSettings(),
		settings.Pagination.Start,
		settings.Pagination.Size)
}

// ParseStarredFilter parses the starred filter to a string.
func (settings *QuerySettings) ParseStarredFilter() string {
	if settings.StarredOnly {
		return `{ "term": { "isStarred": true } }`
	}
	return ""
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
