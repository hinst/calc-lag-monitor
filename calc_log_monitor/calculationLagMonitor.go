package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CalculationLogMonitor struct {
	IsRunning bool
}

func (monitor *CalculationLogMonitor) Start() {
	if monitor.IsRunning {
		panic("Already running")
	}
	monitor.IsRunning = true
	go monitor.run()
}

func (monitor *CalculationLogMonitor) run() {
	for monitor.IsRunning {
		monitor.readCalculationLag()
		time.Sleep(36 * 1000)
	}
}

func (monitor *CalculationLogMonitor) readCalculationLag() {
	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mongo.Connect(context, options.Client().ApplyURI(""))
}
