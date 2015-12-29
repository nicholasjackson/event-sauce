package global

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/facebookgo/inject"
)

type ConfigStruct struct {
	StatsDServerIP string    `json:"stats_d_server_url"`
	RootFolder     string    `json:"root_folder"`
	Data           DataStore `json:"data_store"`
	Queue          Queue     `json:"queue"`
}

type DataStore struct {
	ConnectionString string `json:"connection_string"`
	DataBaseName     string `json:"database_name"`
}

type Queue struct {
	ConnectionString string `json:"connection_string"`
	MessageQueue     string `json:"message_queue"`
}

var Config ConfigStruct

func LoadConfig(config string, rootfolder string) error {
	fmt.Println("Loading Config: ", config)

	file, err := os.Open(config)
	if err != nil {
		return fmt.Errorf("Unable to open config")
	}

	decoder := json.NewDecoder(file)
	Config = ConfigStruct{}
	err = decoder.Decode(&Config)
	Config.RootFolder = rootfolder

	fmt.Println(Config)

	return nil
}

func SetupInjection(objects ...*inject.Object) error {
	var g inject.Graph

	var err error

	err = g.Provide(objects...)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Here the Populate call is creating instances of NameAPI &
	// PlanetAPI, and setting the HTTPTransport on both to the
	// http.DefaultTransport provided above:
	if err := g.Populate(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
