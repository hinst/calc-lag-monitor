package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
)

type App struct {
	Configuration Configuration
	Web           Web
	Storage       *DataStorage
	Monitor       *CalculationLogMonitor
	Provider      *DataProvider
	Exiting       chan bool
}

func (app *App) Run() {
	app.Configuration = LoadConfiguration()
	app.Web.Username = "admin"
	app.Web.Password = app.Configuration.Password
	app.InitializeStorage()
	log.Println("Sampling enabled: " + strconv.FormatBool(app.Configuration.SamplingEnabled))
	if app.Configuration.SamplingEnabled {
		app.InitializeMonitor()
	}
	app.InitializeProvider()
	app.InitializeWebUi()
	app.StartWebServer()
	if app.Monitor != nil {
		app.Monitor.Start()
	}
	if app.Exiting == nil {
		app.Exiting = make(chan bool)
	}
	InstallShutdownReceiver(app.Shutdown)
	<-app.Exiting
	if app.Monitor != nil {
		app.Monitor.Stop()
		app.Monitor.Wait()
	}
	app.Storage.Close()
	log.Print("Shutdown process is now complete")
}

func (app *App) InitializeStorage() {
	if app.Storage == nil {
		app.Storage = &DataStorage{Configuration: &app.Configuration}
		app.Storage.Open()
	}
}

func (app *App) InitializeMonitor() {
	if app.Monitor == nil {
		app.Monitor = &CalculationLogMonitor{
			Configuration: &app.Configuration,
			Storage:       app.Storage,
			LogEnabled:    true,
		}
	}
}

func (app *App) InitializeProvider() {
	if app.Provider == nil {
		app.Provider = &DataProvider{Storage: app.Storage, Configuration: &app.Configuration, Web: &app.Web}
		app.Provider.Register()
	}
}

func (app *App) InitializeWebUi() {
	registerFolder := func(path string) {
		files := http.FileServer(http.Dir("../calc-lag-mon-ui/dist" + path))
		app.Web.Handle(WEB_URL+path, http.StripPrefix(WEB_URL+path, files))
	}
	registerFolder("")
	registerFolder("/_nuxt/")
	registerFolder("/calculation-lag-chart/")
}

func (app *App) StartWebServer() {
	listener, listenerError := net.Listen("tcp", ":3006")
	AssertWrapped(listenerError, "Unable to listen")
	go func() {
		error := http.Serve(listener, nil)
		AssertWrapped(error, "Unable to serve")
	}()
}

func (app *App) Shutdown() {
	log.Print("Received shutdown signal")
	app.Exiting <- true
}
