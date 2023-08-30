package main

import "net/http"

const WEB_URL = "/clm-ui"
const API_URL = "/clm"

type WebFunction = func(http.ResponseWriter, *http.Request)
type WebHandlerWrapper = func(http.ResponseWriter, *http.Request, http.Handler)

type Web struct {
	Username string
	Password string
}

const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_TEXT = "text/plain"

func (web *Web) HandleFunc(path string, f WebFunction) {
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		if web.authenticate(writer, request) {
			f(writer, request)
		}
	})
}

func (web *Web) Handle(path string, handler http.Handler) {
	http.Handle(path, &HandlerWrapper{
		Wrapper: func(writer http.ResponseWriter, request *http.Request, handler http.Handler) {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			if web.authenticate(writer, request) {
				handler.ServeHTTP(writer, request)
			}
		},
		InnerHandler: handler,
	})
}

func (web *Web) authenticate(writer http.ResponseWriter, request *http.Request) (allowed bool) {
	var needCheck = len(web.Password) > 0
	if !needCheck {
		return true
	}
	var username, password, ok = request.BasicAuth()
	allowed = ok && username == web.Username && password == web.Password
	if !allowed {
		writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		writer.WriteHeader(http.StatusUnauthorized)
	}
	return
}

func AddContentTypeHeader(header http.Header, contentType string) {
	header.Add("Content-Type", contentType)
}

type HandlerWrapper struct {
	Wrapper      WebHandlerWrapper
	InnerHandler http.Handler
}

func (me *HandlerWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	me.Wrapper(writer, request, me.InnerHandler)
}
