package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/amoralesc/indexer/zinc"
	"github.com/go-chi/render"
)

// paginateAndSort is a middleware that loads the zinc.QuerySettings
// from the query parameters and adds them as context values.
func paginateAndSort(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the page and page size
		page := r.URL.Query().Get("page")
		pageSize := r.URL.Query().Get("pageSize")
		sortBy := r.URL.Query().Get("sortBy")

		if page == "" {
			page = "1"
		}
		if pageSize == "" {
			pageSize = "0"
		}

		// cast page and page size to int
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("page should be an integer")))
			return
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("pageSize should be an integer")))
			return
		}

		// create the query settings
		querySettings, err := zinc.NewQuerySettings(sortBy, pageInt, pageSizeInt)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// add the query settings to the context
		ctx := context.WithValue(r.Context(), "querySettings", querySettings)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
