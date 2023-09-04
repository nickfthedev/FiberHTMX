package lib

import (
	"io"
	"log"
	"os"
	"testing"
)

// Write Test
func TestLoadConfig(t *testing.T) {
	f, err := os.OpenFile("../tmp/config_test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	filename := "../tmp/config.test.json"
	os.Remove(filename)
	log.Println(filename, " has been removed")
	LoadConfig(filename)
	log.Println(filename, " has been created")
	LoadConfig(filename)

}

func TestWriteConfig(t *testing.T) {
	f, err := os.OpenFile("../tmp/config_test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	c := Config{
		DbDriver: "TESTDRIVER",
	}

	WriteConfig(c, "../tmp/config_test_write.json")
	LoadConfig("../tmp/config_test_write.json")
	if Cfg.DbDriver != "TESTDRIVER" {
		log.Fatalln("Cfg.DbDriver should be TESTDRIVER but is ", Cfg.DbDriver)
	} else {
		log.Println("Passed! Cfg.DbDriver should be TESTDRIVER and is ", Cfg.DbDriver)
	}
}