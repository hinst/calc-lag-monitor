package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	MongoDbUrl              string `json:"mongoDbUrl"`
	MongoDbName             string `json:"mongoDbName"`
	BoltDbFilePath          string `json:"boltDbFilePath"`
	SamplingIntervalSeconds int    `json:"samplingIntervalSeconds"`
	SamplingEnabled         bool   `json:"samplingEnabled"`
}

const configurationFilePath = "./configuration.json"

func LoadConfiguration() Configuration {
	fileContent, fileReadError := os.ReadFile(configurationFilePath)
	AssertWrapped(fileReadError, "Cannot read URL from file "+configurationFilePath)
	var configuration Configuration
	configuration.SetDefault()
	var unmarshalError = json.Unmarshal(fileContent, &configuration)
	AssertWrapped(unmarshalError, "Cannot decode file content from "+configurationFilePath)

	var mongoDbUrl = os.Getenv("MONGO_DB_URL")
	if len(mongoDbUrl) > 0 {
		configuration.MongoDbUrl = mongoDbUrl
	}
	return configuration
}

func (configuration *Configuration) SetDefault() {
	configuration.BoltDbFilePath = "./data.db"
	configuration.SamplingIntervalSeconds = 60
}
