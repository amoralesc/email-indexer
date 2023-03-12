package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/amoralesc/indexer/email"
	"github.com/amoralesc/indexer/utils"
)

var NumParserWorkers, _ = strconv.Atoi(utils.GetenvOrDefault("NUM_PARSER_WORKERS", "8"))
var NumUploaderWorkers, _ = strconv.Atoi(utils.GetenvOrDefault("NUM_UPLOADER_WORKERS", "4"))
var BulkUploadSize, _ = strconv.Atoi(utils.GetenvOrDefault("BULK_UPLOAD_SIZE", "5000"))
var ZincUrl = fmt.Sprintf("http://localhost:%v", utils.GetenvOrDefault("ZINC_PORT", "4080"))
var ZincAdminUser = utils.GetenvOrDefault("ZINC_ADMIN_USER", "admin")
var ZincAdminPassword = utils.GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123")

func parseEmailFiles(files <-chan string, emails chan<- email.Email) {
	for file := range files {
		emailStruct, err := email.EmailFromFile(file)
		if err != nil {
			log.Printf("ERROR: failed to parse %v: %v", file, err)
		} else {
			emails <- emailStruct
		}
	}
}

func uploadEmails(emails <-chan email.Email) {
	bulk := email.BulkEmails{
		Index:   "emails",
		Records: make([]email.Email, BulkUploadSize),
	}
	parsed := 0
	total := 0
	for emailStruct := range emails {
		bulk.Records[parsed] = emailStruct
		parsed++
		if parsed == BulkUploadSize {
			log.Printf("TRACE: uploading %d emails\n", parsed)
			err := email.UploadEmails(bulk, ZincUrl, ZincAdminUser, ZincAdminPassword)
			if err != nil {
				log.Fatal("FATAL: failed to upload emails: ", err)
			}
			total += parsed
			parsed = 0
		}
	}
	if parsed > 0 {
		bulk.Records = bulk.Records[:parsed]
		err := email.UploadEmails(bulk, ZincUrl, ZincAdminUser, ZincAdminPassword)
		if err != nil {
			log.Fatal("FATAL: failed to upload emails: ", err)
		}
		total += parsed
	}
	log.Printf("INFO: goroutine uploaded %d emails\n", total)
}

func parseAndUploadEmails(dir *string) {
	// create channels for passing data between goroutines
	files := make(chan string)
	emails := make(chan email.Email)

	// spawn uploader goroutines
	log.Printf("INFO: spawning %d uploader goroutines", NumUploaderWorkers)
	var wgUploaders sync.WaitGroup
	for i := 0; i < NumUploaderWorkers; i++ {
		wgUploaders.Add(1)
		go func() {
			defer wgUploaders.Done()
			uploadEmails(emails)
		}()
	}

	// spawn file parser goroutines
	log.Printf("INFO: spawning %d parser goroutines", NumParserWorkers)
	var wgParsers sync.WaitGroup
	for i := 0; i < NumParserWorkers; i++ {
		wgParsers.Add(1)
		go func() {
			defer wgParsers.Done()
			parseEmailFiles(files, emails)
		}()
	}

	start := time.Now()
	log.Printf("INFO: starting to upload emails")

	// walk directory and send file paths to channel
	err := filepath.WalkDir(*dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			files <- path
		}
		return nil
	})
	if err != nil {
		log.Fatal("FATAl: failed to walk directory: ", err)
	}

	// close files channel to signal end of parsing
	close(files)
	wgParsers.Wait()
	// close emails channel to signal end of uploading
	close(emails)
	wgUploaders.Wait()

	log.Printf("INFO: finished uploading in %v\n", time.Since(start))
}

func main() {
	// command line flags
	dir := flag.String("d", "", "Directory of the emails. If none is provided, the server will use already indexed emails.")
	flag.Parse()

	if *dir != "" {
		err := email.DeleteIndex(ZincUrl, ZincAdminUser, ZincAdminPassword)
		if err != nil {
			log.Println("WARNING: failed to delete emails index: ", err)
		}
		parseAndUploadEmails(dir)
	}
	// start server
	log.Fatal("FATAL: server not implemented yet")
}
