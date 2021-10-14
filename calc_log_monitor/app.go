package main

import "log"

type App struct {
	Monitor  *CalculationLogMonitor
	Storage  *DataStorage
	Finished chan bool
}

func (app *App) Run() {
	if app.Monitor == nil {
		app.Monitor = &CalculationLogMonitor{
			Configuration: LoadConfiguration(),
		}
	}
	if app.Storage == nil {
		app.Storage = &DataStorage{}
		app.Storage.Open()
	}
	if app.Finished == nil {
		app.Finished = make(chan bool)
	}
	app.Monitor.Start()
	InstallShutdownReceiver(app.Shutdown)
	<-app.Finished
}

func (app *App) Shutdown() {
	log.Print("Received shutdown signal")
	app.Monitor.Stop()
	app.Monitor.Wait()
	log.Print("Shutdown process is now complete")
	app.Finished <- true
}
