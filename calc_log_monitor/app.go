package main

import "log"

type App struct {
	monitor  *CalculationLogMonitor
	finished chan bool
}

func (app *App) Run() {
	if app.finished == nil {
		app.finished = make(chan bool)
	}
	if app.monitor == nil {
		app.monitor = &CalculationLogMonitor{
			Configuration: LoadConfiguration(),
		}
	}
	app.monitor.Start()
	InstallShutdownReceiver(app.Shutdown)
	<-app.finished
}

func (app *App) Shutdown() {
	log.Print("Received shutdown signal")
	app.monitor.Stop()
	app.monitor.Wait()
	log.Print("Shutdown process is now complete")
	app.finished <- true
}
