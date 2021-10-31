package main

import "net/http"

const BASE_URL = "/clm"

type WebFunction = func(http.ResponseWriter, *http.Request)

func HandleFunc(path string, f WebFunction) {
	http.HandleFunc(BASE_URL+path, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		f(writer, request)
	})
}

func AddJsonHeader(header http.Header) {
	header.Add("Content-Type", "application/json")
}
