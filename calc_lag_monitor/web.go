package main

import "net/http"

const BASE_URL = "/clm"

type WebFunction = func(http.ResponseWriter, *http.Request)

func HandleFunc(path string, f WebFunction) {
	http.HandleFunc(BASE_URL+path, f)
}
