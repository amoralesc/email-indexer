package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/amoralesc/email-indexer/indexer/router"
	"github.com/amoralesc/email-indexer/indexer/routines"
	"github.com/amoralesc/email-indexer/indexer/utils"
	"github.com/amoralesc/email-indexer/indexer/zinc"
)

func main() {
	// command line flags
	dir := flag.String("d", "", "Directory of the emails. If none is provided, the server will use already indexed emails.")
	remove := flag.Bool("r", false, "Remove the index before starting the server.")
	flag.Parse()

	port := utils.GetenvOrDefault("PORT", "8080")
	zinc.Service = zinc.NewZincService(fmt.Sprintf("http://localhost:%v", utils.GetenvOrDefault("ZINC_PORT", "4080")), utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin"), utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123"))
	numUploaderWorkers, _ := strconv.Atoi(utils.GetenvOrDefault("NUM_UPLOADER_WORKERS", "4"))
	numParserWorkers, _ := strconv.Atoi(utils.GetenvOrDefault("NUM_PARSER_WORKERS", "8"))
	bulkUploadSize, _ := strconv.Atoi(utils.GetenvOrDefault("BULK_UPLOAD_SIZE", "5000"))

	// remove index if requested
	if *remove {
		log.Printf("INFO: deleting emails index (if exists)")
		err := zinc.Service.DeleteIndex()
		if err != nil {
			log.Println("WARNING: failed to delete emails index: ", err)
		}
	}

	// only parse and upload emails if a directory is provided
	if *dir != "" {
		log.Printf("INFO: creating emails index")
		err := zinc.Service.CreateIndex()
		if err != nil {
			log.Fatal("FATAL: failed to create emails index: ", err)
		}

		log.Println("INFO: starting to parse and upload emails")
		start := time.Now()
		routines.ParseAndUploadEmails(dir, numUploaderWorkers, numParserWorkers, bulkUploadSize)
		log.Printf("INFO: finished uploading in %v\n", time.Since(start))
	}

	log.Println("INFO: starting HTTP server on port ", port)
	r := router.NewRouter()
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
