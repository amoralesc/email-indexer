package routines

import (
	"io/fs"
	"log"
	"path/filepath"
	"sync"

	"github.com/amoralesc/email-indexer/indexer/email"
	"github.com/amoralesc/email-indexer/indexer/zinc"
)

// parseEmailFiles is a routine that parses emails from a channel of file paths
// and sends them to a channel of emails.
func parseEmailFiles(files <-chan string, emails chan<- *email.Email) {
	for file := range files {
		emailObj, err := email.EmailFromFile(file)
		if err != nil {
			log.Printf("WARN: failed to parse %v: %v", file, err)
		} else {
			emails <- emailObj
		}
	}
}

// uploadEmails is a routine that uploads emails from a channel of emails to zinc.
func uploadEmails(emails <-chan *email.Email, bulkUploadSize int, zincAuth *zinc.ZincAuth) {
	bulk := &zinc.BulkEmails{
		Index:   "emails",
		Records: make([]email.Email, bulkUploadSize),
	}
	parsed := 0
	total := 0
	// upload emails in batches of bulkUploadSize
	for emailObj := range emails {
		bulk.Records[parsed] = *emailObj
		parsed++
		if parsed == bulkUploadSize {
			log.Printf("TRACE: uploading %d emails\n", parsed)
			err := zinc.UploadEmails(bulk, zincAuth)
			if err != nil {
				log.Fatal("FATAL: failed to upload emails: ", err)
			}
			total += parsed
			parsed = 0
		}
	}
	if parsed > 0 {
		bulk.Records = bulk.Records[:parsed]
		err := zinc.UploadEmails(bulk, zincAuth)
		if err != nil {
			log.Fatal("FATAL: failed to upload emails: ", err)
		}
		total += parsed
	}
	log.Printf("INFO: goroutine uploaded %d emails, exitting\n", total)
}

// ParseAndUploadEmails is the goroutine manager. It spawns a number of
// goroutines to parse emails from files and upload them to zinc.
func ParseAndUploadEmails(dir string, numUploaderWorkers int, numParserWorkers int, bulkUploadSize int, zincAuth *zinc.ZincAuth) {
	// create channels for passing data between goroutines
	files := make(chan string)
	emails := make(chan *email.Email)

	// spawn uploader goroutines
	log.Printf("TRACE: spawning %d uploader goroutines", numUploaderWorkers)
	var wgUploaders sync.WaitGroup
	for i := 0; i < numUploaderWorkers; i++ {
		wgUploaders.Add(1)
		go func() {
			defer wgUploaders.Done()
			uploadEmails(emails, bulkUploadSize, zincAuth)
		}()
	}

	// spawn file parser goroutines
	log.Printf("TRACE: spawning %d parser goroutines", numParserWorkers)
	var wgParsers sync.WaitGroup
	for i := 0; i < numParserWorkers; i++ {
		wgParsers.Add(1)
		go func() {
			defer wgParsers.Done()
			parseEmailFiles(files, emails)
		}()
	}

	// walk directory and send file paths to channel
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			files <- path
		}
		return nil
	})
	if err != nil {
		log.Fatal("FATAL: failed to walk directory: ", err)
	}

	// close files channel to signal end of parsing
	close(files)
	wgParsers.Wait()
	// close emails channel to signal end of uploading
	close(emails)
	wgUploaders.Wait()
}
