package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	MongoDbUrl  string `json:"mongoDbUrl"`
	MongoDbName string `json:"mongoDbName"`
}

const configurationFilePath = "./configuration.json"

func LoadConfiguration() Configuration {
	fileContent, fileReadError := ioutil.ReadFile(configurationFilePath)
	AssertWrapped(fileReadError, "Cannot read URL from file "+configurationFilePath)
	var configuration Configuration
	var unmarshalError = json.Unmarshal(fileContent, &configuration)
	AssertWrapped(unmarshalError, "Cannot decode file content from "+configurationFilePath)
	return configuration
}
