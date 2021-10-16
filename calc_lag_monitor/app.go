package main

import "log"

type App struct {
	Storage *DataStorage
	Monitor *CalculationLogMonitor
	Exiting chan bool
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
		}
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
