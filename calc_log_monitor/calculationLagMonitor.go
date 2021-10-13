package main

import (
	"context"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tickDuration = 100 * time.Millisecond
const expensiveSourceCount = 10
const trailingCalculationRequestsRatio = 0.1

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

	oldestCalculationRequest := monitor.findCalculationRequests(client, context,
		bson.M{"countOfSources": bson.M{"$lt": expensiveSourceCount}})
	oldestExpensiveCalculationRequest := monitor.findCalculationRequests(client, context,
		bson.M{"countOfSources": bson.M{"$gte": expensiveSourceCount}})
	log.Print(len(oldestCalculationRequest), len(oldestExpensiveCalculationRequest))
}

func (monitor *CalculationLogMonitor) findCalculationRequests(
	client *mongo.Client, context context.Context, query bson.M,
) []CalculationRequest {
	var calculationRequests []CalculationRequest
	calculationRequestCollection := client.
		Database(monitor.Configuration.MongoDbName).
		Collection("calculationRequest")

	countOptions := options.CountOptions{}
	documentCount, documentCountError :=
		calculationRequestCollection.CountDocuments(context, query, &countOptions)
	AssertWrapped(documentCountError, "Unable to read count of calculation requests")

	findOptions := options.FindOptions{}
	findOptions.SetSort(bson.M{"calculateAt": 1})
	var limit = math.Round(float64(documentCount) * trailingCalculationRequestsRatio)
	if limit <= 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	findOptions.SetLimit(int64(limit))
	findOptions.SetProjection(bson.M{"calculateAt": 1})
	cursor, findError := calculationRequestCollection.Find(context, query, &findOptions)
	AssertWrapped(findError, "Unable to find calculation requests")
	defer cursor.Close(context)
	for cursor.Next(context) {
		var calculationRequest CalculationRequest
		cursor.Decode(&calculationRequest)
		decodingError := cursor.Decode(&calculationRequest)
		AssertWrapped(decodingError, "Unable to decode calculation request")
		calculationRequests = append(calculationRequests, calculationRequest)
	}
	AssertWrapped(cursor.Err(), "A cursor error occurred")
	return calculationRequests
}

func (monitor *CalculationLogMonitor) aggregateCalculateAt(calculationRequests []CalculationRequest) (
	min time.Time, average time.Time, max time.Time,
) {
	var haveMin = false
	var haveMax = false
	var sum float64 = 0
	var count = 0
	for _, calculationRequest := range calculationRequests {
		if !haveMin {
			min = calculationRequest.CalculateAt
			haveMin = true
		}
		if calculationRequest.CalculateAt.Before(min) {
			min = calculationRequest.CalculateAt
		}
		if !haveMax {
			max = calculationRequest.CalculateAt
			haveMax = true
		}
		if max.Before(calculationRequest.CalculateAt) {
			max = calculationRequest.CalculateAt
		}
		sum += float64(calculationRequest.CalculateAt.UnixMilli())
		count += 1
	}
	if count > 0 {
		average = time.UnixMilli(int64(sum / float64(count)))
	}
	return
}
