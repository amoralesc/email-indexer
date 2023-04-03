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
	dir := flag.String("d", "", "Directory of the emails. If none is provided, the server will use already indexed emails.")
	remove := flag.Bool("r", false, "Remove the index before starting the server.")
	profiling := flag.Bool("p", false, "Start the profiling server.")
	flag.Parse()

	port := utils.GetenvOrDefault("PORT", "3000")
	profilingPort := utils.GetenvOrDefault("PROFILING_PORT", "6060")
	zinc.Service = zinc.NewZincService(fmt.Sprintf("http://%v:%v", utils.GetenvOrDefault("ZINC_HOST", "localhost"), utils.GetenvOrDefault("ZINC_PORT", "4080")), utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin"), utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123"))
	numUploaderWorkers, _ := strconv.Atoi(utils.GetenvOrDefault("NUM_UPLOADER_WORKERS", "8"))
	numParserWorkers, _ := strconv.Atoi(utils.GetenvOrDefault("NUM_PARSER_WORKERS", "32"))
	bulkUploadSize, _ := strconv.Atoi(utils.GetenvOrDefault("BULK_UPLOAD_SIZE", "5000"))

	_, err := zinc.Service.CheckIndex()
	if err != nil {
		log.Fatal("FATAL: failed to connect to zinc: ", err)
	}

	// start profiling server on goroutine
	if *profiling {
		go func() {
			log.Println("INFO: starting profiling server on port ", profilingPort)
			log.Println(http.ListenAndServe(fmt.Sprintf(":%v", profilingPort), nil))
		}()
	}

	// remove index if requested
	if *remove {
		log.Println("INFO: deleting emails index (if exists)")
		exists, err := zinc.Service.CheckIndex()
		if err != nil {
			log.Fatal("FATAL: failed to check if emails index exists: ", err)
		}

		if exists {
			err := zinc.Service.DeleteIndex()
			if err != nil {
				log.Println("WARN: failed to delete emails index: ", err)
			}
		}
	}

	// only parse and upload emails if a directory is provided
	if *dir != "" {
		exists, err := zinc.Service.CheckIndex()
		if err != nil {
			log.Fatal("FATAL: failed to check if emails index exists: ", err)
		}

		if !exists {
			log.Printf("INFO: creating emails index")
			err := zinc.Service.CreateIndex()
			if err != nil {
				log.Fatal("FATAL: failed to create emails index: ", err)
			}
		}

		log.Println("INFO: starting to parse and upload emails")
		start := time.Now()
		zincAuth := &zinc.ZincAuth{
			Url:      fmt.Sprintf("http://%v:%v", utils.GetenvOrDefault("ZINC_HOST", "localhost"), utils.GetenvOrDefault("ZINC_PORT", "4080")),
			User:     utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin"),
			Password: utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123"),
		}
		routines.ParseAndUploadEmails(dir, numUploaderWorkers, numParserWorkers, bulkUploadSize, zincAuth)
		log.Printf("INFO: finished uploading in %v\n", time.Since(start))
	}

	log.Println("INFO: starting HTTP server on port ", port)
	r := router.NewRouter()
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
