package main

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	MongoDbUrl              string `json:"mongoDbUrl"`
	MongoDbName             string `json:"mongoDbName"`
	SamplingIntervalSeconds int    `json:"samplingIntervalSeconds"`
	SamplingEnabled         bool   `json:"samplingEnabled"`
}

const configurationFilePath = "./configuration.json"

func LoadConfiguration() Configuration {
	fileContent, fileReadError := ioutil.ReadFile(configurationFilePath)
	AssertWrapped(fileReadError, "Cannot read URL from file "+configurationFilePath)
	var configuration Configuration
	configuration.SetDefault()
	var unmarshalError = json.Unmarshal(fileContent, &configuration)
	AssertWrapped(unmarshalError, "Cannot decode file content from "+configurationFilePath)
	return configuration
}

func (configuration *Configuration) SetDefault() {
	configuration.SamplingIntervalSeconds = 60
}
