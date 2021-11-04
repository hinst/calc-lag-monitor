package main

import (
	"flag"
	"runtime/debug"
)

func main() {
	println("STARTING...")
	debug.SetGCPercent(50)
	importEnabled := flag.Bool("import", false, "Import database file")
	db := flag.String("db", "", "File path")
	flag.Parse()
	if *importEnabled {
		var appImporter AppImporter
		appImporter.SourceDbFilePath = *db
		appImporter.Run()
	} else {
		var app App
		app.Run()
	}
	println("...EXITING")
}
