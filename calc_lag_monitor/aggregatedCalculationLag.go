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

func (lag *AggregatedCalculationLagEx) InitializeAggregation(other AggregatedCalculationLag) {
	lag.Min = other.Min
	lag.Max = other.Max
	lag.Count = 0
	lag.Sum = 0
}

func (lag *AggregatedCalculationLagEx) Aggregate(other AggregatedCalculationLag) {
	if other.Min < lag.Min {
		lag.Min = other.Min
	}
	lag.Sum += float64(other.Average)
	if lag.Max < other.Max {
		lag.Max = other.Max
	}
	lag.Count += 1
}

func (lag *AggregatedCalculationLagEx) FinalizeAggregation() {
	if lag.Count > 0 {
		lag.Average = time.Duration(
			int64(lag.Sum / float64(lag.Count)),
		)
	}
}
