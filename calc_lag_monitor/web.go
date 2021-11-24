package main

import "net/http"

const BASE_URL = "/clm"

type WebFunction = func(http.ResponseWriter, *http.Request)

const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_TEXT = "text/plain"

func HandleFunc(path string, f WebFunction) {
	http.HandleFunc(BASE_URL+path, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		f(writer, request)
	})
}

func AddContentTypeHeader(header http.Header, contentType string) {
	header.Add("Content-Type", contentType)
}
