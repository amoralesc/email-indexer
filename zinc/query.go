package zinc

import (
	"fmt"
	"strings"
	"time"
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
