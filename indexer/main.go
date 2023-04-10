package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/amoralesc/email-indexer/indexer/router"
	"github.com/amoralesc/email-indexer/indexer/routines"
	"github.com/amoralesc/email-indexer/indexer/utils"
	"github.com/amoralesc/email-indexer/indexer/zinc"
)

func main() {
	// command line flags
	index := flag.Bool("i", false, "Index the files in the emails directory (env EMAILS_DIR) to zinc.")
	profiling := flag.Bool("p", false, "Start the profiling server.")
	server := flag.Bool("s", false, "Start the emails server (REST API).")
	flag.Parse()

	// env vars
	emailsDir := utils.GetenvOrDefault("EMAILS_DIR", "emails")
	removeIndex, _ := strconv.ParseBool(utils.GetenvOrDefault("REMOVE_INDEX", "false"))
	preventUploadIfIndexExists, _ := strconv.ParseBool(utils.GetenvOrDefault("PREVENT_UPLOAD_IF_INDEX_EXISTS", "false"))

	zinc.StartZincService(fmt.Sprintf("http://%v:%v", utils.GetenvOrDefault("ZINC_HOST", "localhost"), utils.GetenvOrDefault("ZINC_PORT", "4080")), utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin"), utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123"))

	indexExists, err := zinc.Service.CheckIndex()
	if err != nil {
		log.Fatal("FATAL: failed to connect to zinc: ", err)
	}

	// start profiling server on goroutine
	if *profiling {
		profilingPort := utils.GetenvOrDefault("PROFILING_PORT", "6060")
		go func() {
			log.Println("INFO: starting profiling server on port", profilingPort)
			log.Println(http.ListenAndServe(fmt.Sprintf(":%v", profilingPort), nil))
		}()
	}

	// remove index if requested
	if removeIndex {
		if indexExists {
			log.Println("INFO: deleting emails index")
			err := zinc.Service.DeleteIndex()
			if err != nil {
				log.Panic("ERROR: failed to delete emails index:", err)
			}
		}
	}

	// index the emails
	if *index {
		// recheck index, might have been deleted
		indexExists, err := zinc.Service.CheckIndex()
		if err != nil {
			log.Fatal("FATAL: failed to check if emails index exists: ", err)
		}

		if preventUploadIfIndexExists && indexExists {
			log.Println("INFO: emails index already exists, skipping upload")
		} else {
			if !indexExists {
				log.Printf("INFO: creating emails index")
				err := zinc.Service.CreateIndex()
				if err != nil {
					log.Fatal("FATAL: failed to create emails index: ", err)
				}
			}

			log.Println("INFO: starting to parse and upload emails at dir:", emailsDir)
			start := time.Now()
			zincAuth := &zinc.ZincAuth{
				Url:      fmt.Sprintf("http://%v:%v", utils.GetenvOrDefault("ZINC_HOST", "localhost"), utils.GetenvOrDefault("ZINC_PORT", "4080")),
				User:     utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin"),
				Password: utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123"),
			}
			numUploaderWorkers, _ := strconv.Atoi(utils.GetenvOrDefault("NUM_UPLOADER_WORKERS", "32"))
			numParserWorkers, _ := strconv.Atoi(utils.GetenvOrDefault("NUM_PARSER_WORKERS", "128"))
			bulkUploadSize, _ := strconv.Atoi(utils.GetenvOrDefault("BULK_UPLOAD_SIZE", "5000"))
			routines.ParseAndUploadEmails(emailsDir, numUploaderWorkers, numParserWorkers, bulkUploadSize, zincAuth)
			log.Printf("INFO: finished uploading in %v\n", time.Since(start))
		}
	}

	if !*server {
		return
	}
	port := utils.GetenvOrDefault("API_PORT", "3000")
	log.Println("INFO: starting REST API on port", port)
	r := router.NewRouter()
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
