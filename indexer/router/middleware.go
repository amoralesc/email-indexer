package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/amoralesc/email-indexer/indexer/zinc"
	"github.com/go-chi/render"
)

// loadQuerySettings is a middleware that loads the zinc.QuerySettings
// from the query parameters and adds them as context values.
func loadQuerySettings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the page and page size
		start := r.URL.Query().Get("start")
		size := r.URL.Query().Get("size")
		sortBy := r.URL.Query().Get("sortBy")
		starredOnly := r.URL.Query().Get("starredOnly")

		if start == "" {
			start = "1"
		}
		if size == "" {
			size = "0"
		}
		if starredOnly == "" {
			starredOnly = "false"
		}

		// cast start and size to int
		startInt, err := strconv.Atoi(start)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("start should be an integer")))
			return
		}
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("size should be an integer")))
			return
		}
		// cast starredOnly to bool
		starredOnlyBool, err := strconv.ParseBool(starredOnly)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("starredOnly should be a boolean")))
			return
		}

		// create the query settings
		querySettings, err := zinc.NewQuerySettings(sortBy, startInt, sizeInt, starredOnlyBool)
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
