package main

import "runtime/debug"

func main() {
	debug.SetGCPercent(50)
	println("STARTING...")
	var app App
	app.Run()
	println("...EXITING")
}
