package main

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CalculationRangeMonitor struct {
	Storage       *DataStorage
	Configuration *Configuration

	url       string
	isRunning bool
	ticker    Ticker
	finished  chan bool
}

func (monitor *CalculationRangeMonitor) Start() {
	interval := time.Second * time.Duration(monitor.Configuration.SamplingIntervalSeconds)
	monitor.ticker.Initialize(interval)
	monitor.url = monitor.Configuration.MongoDbUrl
}

func (monitor *CalculationRangeMonitor) Stop() {
	monitor.isRunning = false
}

func (monitor *CalculationRangeMonitor) Wait() {
	<-monitor.finished
}

func (monitor *CalculationRangeMonitor) run() {
	for monitor.isRunning {
		if monitor.ticker.Advance(tickDuration) {
			monitor.runOnceSafe()
		}
		time.Sleep(tickDuration)
	}
	monitor.finished <- true
}

func (monitor *CalculationRangeMonitor) runOnceSafe() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Print("Unable to process calculation log\n", recovered, "\n", string(debug.Stack()))
		}
	}()
	monitor.runOnce()
}

func (monitor *CalculationRangeMonitor) runOnce() {
}

func (monitor *CalculationRangeMonitor) readCalculationRangeInfo() (result CalculationRangeInfo) {
	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, mongoError := mongo.Connect(context, options.Client().ApplyURI(monitor.url))
	AssertWrapped(mongoError, "Unable to connect to MongoDB at URL "+monitor.url)
	disconnect := func() {
		disconnectError := client.Disconnect(context)
		AssertWrapped(disconnectError, "Unable to disconnect")
	}
	defer disconnect()
	return
}
