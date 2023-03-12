package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/amoralesc/indexer/email"
)

func GetenvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var NumParserWorkers, _ = strconv.Atoi(GetenvOrDefault("NUM_PARSER_WORKERS", "8"))
var NumUploaderWorkers, _ = strconv.Atoi(GetenvOrDefault("NUM_UPLOADER_WORKERS", "4"))
var BulkUploadSize, _ = strconv.Atoi(GetenvOrDefault("BULK_UPLOAD_SIZE", "5000"))
var ZincPort = GetenvOrDefault("ZINC_PORT", "4080")
var ZincUrl = fmt.Sprintf("http://localhost:%v/api/_bulkv2", ZincPort)
var ZincAdminUser = GetenvOrDefault("ZINC_ADMIN_USER", "admin")
var ZincAdminPassword = GetenvOrDefault("ZINC_ADMIN_PASSWORD", "Complexpass#123")

type BulkEmails struct {
	Index   string        `json:"index"`
	Records []email.Email `json:"records"`
}

func UploadEmails(bulk BulkEmails) error {
	// Convert the struct to JSON
	jsonBytes, err := json.Marshal(bulk)
	if err != nil {
		return err
	}
	jsonReader := bytes.NewReader(jsonBytes)

	// Create the request
	req, err := http.NewRequest("POST", ZincUrl, jsonReader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(ZincAdminUser, ZincAdminPassword)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the response
	// log.Println("Zinc responds with status:", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc didn't accept bulk load: %v", string(body))
	}

	return nil
}

func parseEmailFiles(files <-chan string, results chan<- email.Email, errs chan<- error) {
	for file := range files {
		emailStruct, err := email.EmailFromFile(file)
		if err != nil {
			errs <- fmt.Errorf("failed to parse %v: %v", file, err)
		} else {
			results <- emailStruct
		}
	}
}

func uploadEmails(emails <-chan email.Email, errs chan<- error) {
	bulk := BulkEmails{
		Index:   "emails",
		Records: make([]email.Email, BulkUploadSize),
	}
	parsed := 0
	total := 0
	for emailStruct := range emails {
		bulk.Records[parsed] = emailStruct
		parsed++
		if parsed == BulkUploadSize {
			log.Printf("Uploading %d emails\n", parsed)
			err := UploadEmails(bulk)
			if err != nil {
				errs <- fmt.Errorf("error uploading emails: %v", err)
				return
			}
			total += parsed
			parsed = 0
		}
	}
	if parsed > 0 {
		bulk.Records = bulk.Records[:parsed]
		err := UploadEmails(bulk)
		if err != nil {
			errs <- fmt.Errorf("error uploading emails: %v", err)
		}
		total += parsed
	}
	log.Printf("Worker uploaded %d emails\n", total)
}

func main() {
	// Command line flags
	dir := flag.String("d", "", "Directory of the emails")
	// logFailed := flag.Bool("l", false, "Log emails that failed to parse")
	flag.Parse()

	// create channels for passing data and errors between goroutines
	files := make(chan string)
	emails := make(chan email.Email)
	errs := make(chan error)

	// spawn uploader goroutine
	var wg sync.WaitGroup
	for i := 0; i < NumUploaderWorkers; i++ {
		wg.Add(1)
		go func() {
			uploadEmails(emails, errs)
			wg.Done()
		}()
	}

	// spawn file parser goroutines
	var wg2 sync.WaitGroup
	for i := 0; i < NumParserWorkers; i++ {
		wg2.Add(1)
		go func() {
			parseEmailFiles(files, emails, errs)
			wg2.Done()
		}()
	}

	// goroutine to log errors from errs channel
	go func() {
		for err := range errs {
			log.Println(err)
		}
	}()

	start := time.Now()
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
		log.Fatal("error walking directory:", err)
	}

	// close file channel to signal end of parsing
	close(files)

	// wait for all file parser goroutines to finish
	wg2.Wait()

	// close email channel to signal end of uploading
	close(emails)

	// wait for uploader goroutine to finish
	wg.Wait()

	// close error channel to signal end of logging
	close(errs)

	log.Printf("Finished in %v", time.Since(start))
}
