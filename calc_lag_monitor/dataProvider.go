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
}

func (provider *DataProvider) Lag(responseWriter http.ResponseWriter, request *http.Request) {
	start := ParseIntOr0(request.URL.Query().Get("start"))
	end := ParseIntOr0(request.URL.Query().Get("end"))
	aggregatedRows := provider.Storage.ReadCalculationLagInfoRows(start, end)
	AddJsonHeader(responseWriter.Header())
	responseWriter.Write(EncodeJson(aggregatedRows))
}

func (provider *DataProvider) DbStats(responseWriter http.ResponseWriter, request *http.Request) {
	AddJsonHeader(responseWriter.Header())
	responseWriter.Write(EncodeJson(provider.Storage.GetStatistics()))
}
