package main

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

type AggregatedCalculationLag struct {
	Min     time.Duration
	Average time.Duration
	Max     time.Duration
}

func (lag *AggregatedCalculationLag) Write(buffer *bytes.Buffer) {
	BinaryWrite(buffer, BinaryObjectVersionNumber(1))
	BinaryWrite(buffer, lag.Min)
	BinaryWrite(buffer, lag.Average)
	BinaryWrite(buffer, lag.Max)
}

func (lag *AggregatedCalculationLag) Read(buffer *bytes.Buffer) {
	var version BinaryObjectVersionNumber
	BinaryRead(buffer, &version)
	if version != 1 {
		panic(errors.New("Expected version 1 but got " + strconv.Itoa(int(version))))
	}
	BinaryRead(buffer, &lag.Min)
	BinaryRead(buffer, &lag.Average)
	BinaryRead(buffer, &lag.Max)
}

func (lag *AggregatedCalculationLag) ReadFromRequest(request *AggregatedCalculationRequest) {
	now := time.Now()
	lag.Min = now.Sub(request.Min)
	lag.Average = now.Sub(request.Average)
	lag.Max = now.Sub(request.Max)
}

type AggregatedCalculationLagEx struct {
	AggregatedCalculationLag
	Count int
	Sum   float64
}

func (lag *AggregatedCalculationLagEx) InitializeAggregation(firstItem AggregatedCalculationLag) {
	lag.Min = firstItem.Min
	lag.Max = firstItem.Max
	lag.Count = 0
	lag.Sum = 0
}

func (lag *AggregatedCalculationLagEx) Aggregate(item AggregatedCalculationLag) {
	if item.Min < lag.Min {
		lag.Min = item.Min
	}
	lag.Sum += float64(item.Average)
	if lag.Max < item.Max {
		lag.Max = item.Max
	}
	lag.Count += 1
}

func (lag *AggregatedCalculationLagEx) FinalizeAggregation() AggregatedCalculationLag {
	if lag.Count > 0 {
		lag.Average = time.Duration(
			int64(lag.Sum / float64(lag.Count)),
		)
	}
	return lag.AggregatedCalculationLag
}
