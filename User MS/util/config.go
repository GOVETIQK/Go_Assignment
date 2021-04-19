package util

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Host string `json:"host"`
	} `json:"database"`

	Kafka struct {
		Broker string `json:"broker"`
		Topic  string `json:"Topic"`
	} `json:"kafka"`

	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
		panic("Error Loading Config File")
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal(err)
		panic("Invalid Config File")
	}
	return config
}
