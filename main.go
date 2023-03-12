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
	log.Println("Zinc responds with status:", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc didn't accept bulk load: %v", string(body))
	}

	return nil
}

func main() {
	// Command line flags
	dir := flag.String("d", "", "Directory of the emails")
	logFailed := flag.Bool("l", false, "Log emails that failed to parse")
	flag.Parse()

	// Only upload emails if a directory was specified
	// Otherwise, just run the server
	if *dir != "" {
		log.Println("Parsing emails to JSON and uploading...")
		start := time.Now()
		bulk := BulkEmails{
			Index:   "emails",
			Records: make([]email.Email, BulkUploadSize),
		}
		parsed := 0
		total := 0
		// Walk the directory and convert each valid email file to JSON
		err := filepath.WalkDir(*dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				// Upload the emails in batches of BulkUploadSize
				emailStruct, err := email.EmailFromFile(path)
				if err != nil {
					if *logFailed {
						log.Printf("failed to parse %v: %v\n", path, err)
					}
				} else {
					bulk.Records[parsed] = emailStruct
					parsed++
					if parsed == BulkUploadSize {
						log.Printf("Uploading %d emails...\n", parsed)
						err = UploadEmails(bulk)
						if err != nil {
							log.Fatal("error uploading emails: ", err)
						}
						total += parsed
						parsed = 0
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal("error walking directory:", err)
		}
		// Upload the remaining emails
		if parsed > 0 {
			bulk.Records = bulk.Records[:parsed]
			log.Printf("Uploading %d emails...\n", parsed)
			err = UploadEmails(bulk)
			if err != nil {
				log.Fatal("error uploading emails: ", err)
			}
			total += parsed
		}

		log.Println("Done uploading emails in", time.Since(start))
		log.Printf("Uploaded %d emails.\n", total)
	}

}
