package main

import (
	"net/http"
)

type DataProvider struct {
	Storage *DataStorage
}

func (provider *DataProvider) Register() {
	HandleFunc("/lag", provider.Lag)
	HandleFunc("/dbStats", provider.DbStats)
	HandleFunc("/removeAnomalies", provider.RemoveAnomalies)
}

func (provider *DataProvider) Lag(responseWriter http.ResponseWriter, request *http.Request) {
	start := ParseIntOr0(request.URL.Query().Get("start"))
	end := ParseIntOr0(request.URL.Query().Get("end"))
	aggregatedRows := provider.Storage.ReadCalculationLagInfoRows(start, end)
	AddContentTypeHeader(responseWriter.Header(), CONTENT_TYPE_JSON)
	responseWriter.Write(EncodeJson(aggregatedRows))
}

func (provider *DataProvider) DbStats(responseWriter http.ResponseWriter, request *http.Request) {
	AddContentTypeHeader(responseWriter.Header(), CONTENT_TYPE_JSON)
	responseWriter.Write(EncodeJson(provider.Storage.GetStatistics()))
}

func (provider *DataProvider) RemoveAnomalies(responseWriter http.ResponseWriter, request *http.Request) {
	AddContentTypeHeader(responseWriter.Header(), CONTENT_TYPE_TEXT)
	writeEnabled := request.URL.Query().Get("writeEnabled")
	responseWriter.Write([]byte(provider.Storage.RemoveAnomalies(len(writeEnabled) > 0)))
}
