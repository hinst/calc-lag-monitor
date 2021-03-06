package main

import (
	"context"
	"log"
	"math"
	"runtime/debug"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tickDuration = 100 * time.Millisecond
const expensiveSourceCount = 10
const trailingCalculationRequestsRatio = 0.1

type CalculationLogMonitor struct {
	Storage       *DataStorage
	Configuration *Configuration
	LogEnabled    bool

	url       string
	isRunning bool
	ticker    Ticker
	finished  chan bool
}

func (monitor *CalculationLogMonitor) Start() {
	interval := time.Second * time.Duration(monitor.Configuration.SamplingIntervalSeconds)
	monitor.ticker.Initialize(interval)
	if monitor.finished == nil {
		monitor.finished = make(chan bool)
	}
	monitor.url = monitor.Configuration.MongoDbUrl

	if monitor.isRunning {
		panic(CreateException("Already running", nil))
	}
	monitor.isRunning = true
	go monitor.run()
}

func (monitor *CalculationLogMonitor) Stop() {
	monitor.isRunning = false
}

func (monitor *CalculationLogMonitor) Wait() {
	<-monitor.finished
}

func (monitor *CalculationLogMonitor) run() {
	for monitor.isRunning {
		if monitor.ticker.Advance(tickDuration) {
			monitor.runOnceSafe()
		}
		time.Sleep(tickDuration)
	}
	monitor.finished <- true
}

func (monitor *CalculationLogMonitor) runOnceSafe() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Print("Unable to process calculation log\n", recovered, "\n", string(debug.Stack()))
		}
	}()
	monitor.runOnce()
}

func (monitor *CalculationLogMonitor) runOnce() {
	calculationLagInfoRow := monitor.readCalculationLag()
	if monitor.LogEnabled {
		log.Println(calculationLagInfoRow.String())
	}
	monitor.Storage.SaveCalculationLagInfoRow(&calculationLagInfoRow)
}

func (monitor *CalculationLogMonitor) readCalculationLag() CalculationLagInfoRow {
	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, mongoError := mongo.Connect(context, options.Client().ApplyURI(monitor.url))
	AssertWrapped(mongoError, "Unable to connect to MongoDB at URL "+monitor.url)
	disconnect := func() {
		disconnectError := client.Disconnect(context)
		AssertWrapped(disconnectError, "Unable to disconnect")
	}
	defer disconnect()

	oldestCheapCalculationRequests := monitor.findCalculationRequests(client, context,
		bson.M{"countOfSources": bson.M{"$lt": expensiveSourceCount}})
	expensiveOldestCalculationRequests := monitor.findCalculationRequests(client, context,
		bson.M{"countOfSources": bson.M{"$gte": expensiveSourceCount}})
	cheapAggregatedRequest := monitor.aggregateCalculateAt(oldestCheapCalculationRequests)
	expensiveAggregatedRequest := monitor.aggregateCalculateAt(expensiveOldestCalculationRequests)

	var cheapAggregatedLag AggregatedCalculationLag
	if cheapAggregatedRequest != nil {
		cheapAggregatedLag.ReadFromRequest(cheapAggregatedRequest)
	}
	var expensiveAggregatedLag AggregatedCalculationLag
	if expensiveAggregatedRequest != nil {
		expensiveAggregatedLag.ReadFromRequest(expensiveAggregatedRequest)
	}
	calculationLagInfoRow := CalculationLagInfoRow{
		Time:      time.Now(),
		Cheap:     cheapAggregatedLag,
		Expensive: expensiveAggregatedLag}
	return calculationLagInfoRow
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

func (monitor *CalculationLogMonitor) aggregateCalculateAt(calculationRequests []CalculationRequest) *AggregatedCalculationRequest {
	if len(calculationRequests) <= 0 {
		return nil
	}
	var aggregated AggregatedCalculationRequest
	var haveMin = false
	var haveMax = false
	var sum float64 = 0
	var count = 0
	for _, calculationRequest := range calculationRequests {
		if !haveMin {
			aggregated.Min = calculationRequest.CalculateAt
			haveMin = true
		}
		if calculationRequest.CalculateAt.Before(aggregated.Min) {
			aggregated.Min = calculationRequest.CalculateAt
		}
		if !haveMax {
			aggregated.Max = calculationRequest.CalculateAt
			haveMax = true
		}
		if aggregated.Max.Before(calculationRequest.CalculateAt) {
			aggregated.Max = calculationRequest.CalculateAt
		}
		sum += float64(calculationRequest.CalculateAt.UnixMilli())
		count += 1
	}
	if count > 0 {
		aggregated.Average = time.UnixMilli(int64(sum / float64(count)))
	}
	return &aggregated
}
