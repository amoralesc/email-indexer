package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/amoralesc/indexer/routines"
	"github.com/amoralesc/indexer/utils"
	"github.com/amoralesc/indexer/zinc"
)

var NumUploaderWorkers, _ = strconv.Atoi(utils.GetenvOrDefault("NUM_UPLOADER_WORKERS", "4"))
var NumParserWorkers, _ = strconv.Atoi(utils.GetenvOrDefault("NUM_PARSER_WORKERS", "8"))
var BulkUploadSize, _ = strconv.Atoi(utils.GetenvOrDefault("BULK_UPLOAD_SIZE", "5000"))
var ZincUrl = fmt.Sprintf("http://localhost:%v", utils.GetenvOrDefault("ZINC_PORT", "4080"))
var ZincAdminUser = utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin")
var ZincAdminPassword = utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123")

func main() {
	// command line flags
	dir := flag.String("d", "", "Directory of the emails. If none is provided, the server will use already indexed emails.")
	flag.Parse()

	zincServer := zinc.ServerAuth{
		Url:      ZincUrl,
		User:     ZincAdminUser,
		Password: ZincAdminPassword,
	}

	// only parse and upload emails if a directory is provided
	if *dir != "" {
		log.Printf("INFO: deleting emails index (if exists)")
		err := zinc.DeleteIndex("emails", zincServer)
		if err != nil {
			log.Println("WARNING: failed to delete emails index: ", err)
		}

		log.Printf("INFO: creating emails index")
		err = zinc.CreateIndex(zincServer)
		if err != nil {
			log.Fatal("FATAL: failed to create emails index: ", err)
		}

		log.Println("INFO: starting to parse and upload emails")
		start := time.Now()
		routines.ParseAndUploadEmails(dir, NumUploaderWorkers, NumParserWorkers, BulkUploadSize, zincServer)
		log.Printf("INFO: finished uploading in %v\n", time.Since(start))
	}
	// start server
	log.Fatal("FATAL: server not implemented yet")
}
