package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/amoralesc/indexer/zinc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// NewRouter creates a new router
func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/emails", func(r chi.Router) {
		r.With(paginateAndSort).Get("/", ListEmails)
		r.With(paginateAndSort).Get("/search", SearchEmails)
		r.With(paginateAndSort).Get("/query", QueryEmails)

		r.Route("/{emailId}", func(r chi.Router) {
			r.Get("/", GetEmailById)
		})
	})

	r.Get("/api/addresses", ListAddresses)

	return r
}

// ListEmails returns a list of all emails in zinc.
func ListEmails(w http.ResponseWriter, r *http.Request) {
	querySettings := r.Context().Value("querySettings").(*zinc.QuerySettings)
	resp, err := zinc.Service.GetAllEmails(querySettings)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		if strings.Contains(err.Error(), "connection refused") {
			render.Render(w, r, ErrServiceUnavailable)
			return
		}

		render.Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// SearchEmails returns a list of emails that match the search query.
// The search query comes from the body of the request as a JSON object.
func SearchEmails(w http.ResponseWriter, r *http.Request) {
	querySettings := r.Context().Value("querySettings").(*zinc.QuerySettings)
	var searchQuery *zinc.SearchQuery

	// get the search query from the body
	if err := render.DecodeJSON(r.Body, &searchQuery); err != nil {
		log.Printf("ERROR: %v\n", err)
		if err.Error() == "EOF" {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("search query must contain at least one field")))
			return
		}

		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp, err := zinc.Service.GetEmailsBySearchQuery(searchQuery, querySettings)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		if strings.Contains(err.Error(), "connection refused") {
			render.Render(w, r, ErrServiceUnavailable)
			return
		}

		render.Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// QueryEmails returns a list of emails that match the query string.
// The query string comes as a query parameter.
func QueryEmails(w http.ResponseWriter, r *http.Request) {
	querySettings := r.Context().Value("querySettings").(*zinc.QuerySettings)
	queryString := r.URL.Query().Get("q")
	if queryString == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("query string can't be empty")))
		return
	}

	resp, err := zinc.Service.GetEmailsByQueryString(queryString, querySettings)

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		if strings.Contains(err.Error(), "connection refused") {
			render.Render(w, r, ErrServiceUnavailable)
			return
		}

		render.Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// GetEmailById returns an email by its message id.
func GetEmailById(w http.ResponseWriter, r *http.Request) {
	resp, err := zinc.Service.GetEmailByMessageId(chi.URLParam(r, "emailId"))

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		if strings.Contains(err.Error(), "connection refused") {
			render.Render(w, r, ErrServiceUnavailable)
			return
		}

		render.Render(w, r, ErrInternalServer)
		return
	}

	if len(resp.Emails) == 0 {
		render.Render(w, r, ErrNotFound)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// ListAddresses returns a list of all email addresses in zinc.
func ListAddresses(w http.ResponseWriter, r *http.Request) {
	resp, err := zinc.Service.GetAllEmailAddresses()

	if err != nil {
		log.Printf("ERROR: %v\n", err)
		if strings.Contains(err.Error(), "connection refused") {
			render.Render(w, r, ErrServiceUnavailable)
			return
		}

		render.Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
