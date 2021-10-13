package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tickDuration = 100 * time.Millisecond
const expensiveSourceCount = 10

type CalculationLogMonitor struct {
	Configuration Configuration
	IsRunning     bool
	Url           string
	Finished      chan bool
	Interval      time.Duration
	ticker        Ticker
}

func (monitor *CalculationLogMonitor) Start() {
	if monitor.Interval == 0 {
		monitor.Interval = time.Second * 3
	}
	monitor.ticker.Initialize(monitor.Interval)
	if monitor.Finished == nil {
		monitor.Finished = make(chan bool)
	}
	if monitor.IsRunning {
		panic(CreateException("Already running", nil))
	}
	monitor.Url = monitor.Configuration.MongoDbUrl
	monitor.IsRunning = true
	go monitor.run()
}

func (monitor *CalculationLogMonitor) Stop() {
	monitor.IsRunning = false
}

func (monitor *CalculationLogMonitor) Wait() {
	<-monitor.Finished
}

func (monitor *CalculationLogMonitor) run() {
	for monitor.IsRunning {
		if monitor.ticker.Advance(tickDuration) {
			monitor.readCalculationLagSafe()
		}
		time.Sleep(tickDuration)
	}
	monitor.Finished <- true
}

func (monitor *CalculationLogMonitor) readCalculationLagSafe() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Print("Unable to read calculation log\n", recovered)
		}
	}()
	monitor.readCalculationLag()
}

func (monitor *CalculationLogMonitor) readCalculationLag() {
	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, mongoError := mongo.Connect(context, options.Client().ApplyURI(monitor.Url))
	AssertWrapped(mongoError, "Unable to connect to MongoDB at URL "+monitor.Url)
	disconnect := func() {
		disconnectError := client.Disconnect(context)
		AssertWrapped(disconnectError, "Unable to disconnect")
	}
	defer disconnect()

	oldestCalculationRequest := monitor.findCalculationRequest(client, context,
		bson.M{"countOfSources": bson.M{"$lt": expensiveSourceCount}})
	oldestExpensiveCalculationRequest := monitor.findCalculationRequest(client, context,
		bson.M{"countOfSources": bson.M{"$gte": expensiveSourceCount}})
	log.Print(oldestCalculationRequest.CalculateAt, oldestExpensiveCalculationRequest.CalculateAt)
}

func (monitor *CalculationLogMonitor) findCalculationRequest(
	client *mongo.Client, context context.Context, query bson.M,
) CalculationRequest {
	var calculationRequest CalculationRequest
	calculationRequestCollection := client.
		Database(monitor.Configuration.MongoDbName).
		Collection("calculationRequest")
	findOptions := options.FindOneOptions{}
	findOptions.SetSort(bson.M{"calculateAt": 1})
	oldestCalculationRequest := calculationRequestCollection.FindOne(context, query, &findOptions)
	if oldestCalculationRequest != nil {
		AssertWrapped(oldestCalculationRequest.Err(), "Unable to read oldest calculation request")
		decodingError := oldestCalculationRequest.Decode(&calculationRequest)
		AssertWrapped(decodingError, "Unable to decode calculation request")
	}
	return calculationRequest
}
