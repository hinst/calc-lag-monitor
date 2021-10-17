package main

import (
	"log"
	"net"
	"net/http"
)

type App struct {
	Storage  *DataStorage
	Monitor  *CalculationLogMonitor
	Provider *DataProvider
	Exiting  chan bool
}

func (app *App) Run() {
	if app.Storage == nil {
		app.Storage = &DataStorage{}
		app.Storage.Open()
	}
	if app.Monitor == nil {
		app.Monitor = &CalculationLogMonitor{
			Configuration: LoadConfiguration(),
			Storage:       app.Storage,
			LogEnabled:    true,
		}
	}
	if app.Provider == nil {
		app.Provider = &DataProvider{Storage: app.Storage}
		app.Provider.Register()
		listener, listenerError := net.Listen("tcp", ":3006")
		AssertWrapped(listenerError, "Unable to listen")
		go func() {
			error := http.Serve(listener, nil)
			AssertWrapped(error, "Unable to serve")
		}()
	}
	if app.Exiting == nil {
		app.Exiting = make(chan bool)
	}
	app.Monitor.Start()
	InstallShutdownReceiver(app.Shutdown)
	<-app.Exiting
	app.Monitor.Stop()
	app.Monitor.Wait()
	app.Storage.Close()
	log.Print("Shutdown process is now complete")
}

func (app *App) Shutdown() {
	log.Print("Received shutdown signal")
	app.Exiting <- true
}
