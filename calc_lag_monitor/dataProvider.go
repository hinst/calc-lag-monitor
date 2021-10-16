package main

import (
	"log"
	"net/http"
)

type DataProvider struct {
	Storage *DataStorage
}

func (provider *DataProvider) Register() {
	HandleFunc("/lag", provider.Lag)
}

func (provider *DataProvider) Lag(responseWriter http.ResponseWriter, request *http.Request) {
	start := ParseIntOr0(request.URL.Query().Get("start"))
	end := ParseIntOr0(request.URL.Query().Get("end"))
	rows := provider.Storage.ReadCalculationLagInfoRows(start, end)
	log.Print(len(rows))
}
