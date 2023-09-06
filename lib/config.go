package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var Config ConfigStruct

type ConfigStruct struct {
	//
	AppName string
	//Database Driver SQLITE, POSTGRESQL
	DbDriver string
	//Database Filename (for SQLITE)
	DbFileName string
	// For PostgreSQL
	DbHost string
	DbUser string
	DbPass string
	DbName string
	DbPort string

	// Use URL instead of user password server and so on
	PostgresURL string
	// Secret Token (do not share!)
	TokenSecret string
	// API Port
	Port string
	// API URL
	Host string

	// SMTP
	SMTPPort string
	SMTPHost string
	SMTPUser string
	SMTPPass string
}

// }

// Load config.json from root. If the file does not exist, a config.json will be created
func LoadConfig(path string) (ConfigStruct, error) {
	var payload ConfigStruct
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err, "\n=============================\n", "Creating config.json!", "\n=============================")
		//If config does not exist yet, create one
		payload.AppName = "MyApp"

		payload.TokenSecret = "ABCDEFG"

		payload.Port = "3000"
		payload.Host = "127.0.0.1"

		payload.DbDriver = "SQLITE"
		payload.DbFileName = "database.sqlite"

		payload.DbHost = "YOUR HOST IF POSTGRESQL is used"
		payload.DbName = "YOUR DB NAME IF POSTGRESQL is used"
		payload.DbPass = "YOUR PASSWORD IF POSTGRESQL is used"
		payload.DbUser = "DB USERNAME IF POSTGRESQL is used"
		payload.DbPort = "DB PORT IF POSTGRESQL is used"
		payload.PostgresURL = "PREFERRED WAY OF CONNECTING"

		payload.SMTPHost = ""
		payload.SMTPPort = ""
		payload.SMTPUser = ""
		payload.SMTPPass = ""
		WriteConfig(payload, path)

	}

	fmt.Println("Successfully Opened config.json")
	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &payload)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	Config = payload
	return payload, nil
}

// Writes given config struct into config.json at the root directory
func WriteConfig(cfg ConfigStruct, path string) {

	file, _ := json.MarshalIndent(cfg, "", " ")
	_ = os.WriteFile(path, file, 0644)
}
