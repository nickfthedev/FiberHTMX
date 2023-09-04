package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Config ConfigStruct


type ConfigStruct struct {
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
	APIPort string
	// API URL
	APIUrl string
	// List of Plugin structs
	//PluginList []Plugin
}

// type Plugin struct {
// 	Name string
// 	Port string
// 	// Determines if the executable should be build before run.
// 	// Requires golang compiler and source code of plugin in folder
// 	BuildExecutable bool
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
		payload.TokenSecret = "ABCDEFG"

		payload.APIPort = "3000"
		payload.APIUrl = "127.0.0.1"

		payload.DbDriver = "SQLITE"
		payload.DbFileName = "database.sqlite"

		payload.DbHost = "YOUR HOST IF POSTGRESQL is used"
		payload.DbName = "YOUR DB NAME IF POSTGRESQL is used"
		payload.DbPass = "YOUR PASSWORD IF POSTGRESQL is used"
		payload.DbUser = "DB USERNAME IF POSTGRESQL is used"
		payload.DbPort = "DB PORT IF POSTGRESQL is used"
		payload.PostgresURL = "PREFERRED WAY OF CONNECTING"

		//var pluginExampleList []Plugin
		//pluginExampleList = append(pluginExampleList, Plugin{Name: "plugin-example", Port: "3002"})
		//payload.PluginList = pluginExampleList

		WriteConfig(payload, path)

	}

	fmt.Println("Successfully Opened config.json")
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

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
	_ = ioutil.WriteFile(path, file, 0644)
}
