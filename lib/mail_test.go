package lib

import (
	"log"
	"testing"
)

// Write Test
func TestSendMail(t *testing.T) {
	// Load config
	LoadConfig("../config.json") //Config must exist with real login data
	// Send Mail
	err := SendEmail(Config.SMTPUser, "This is a Test Mail", "This is a Test Message", "text/plain")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
