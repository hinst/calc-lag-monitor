package main

import (
	"flag"
	"runtime/debug"
)

func main() {
	println("STARTING...")
	debug.SetGCPercent(50)
	importEnabled := flag.Bool("import", false, "Import database file")
	destination := flag.String("destination", "", "Destination file path")
	source := flag.String("source", "", "Source file path")
	flag.Parse()
	if *importEnabled {
		var appImporter AppImporter
		appImporter.DestinationDbFilePath = *destination
		appImporter.SourceDbFilePath = *source
		appImporter.Run()
	} else {
		var app App
		app.Run()
	}
	println("...EXITING")
}
