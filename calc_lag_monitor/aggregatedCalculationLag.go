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

func (lag *AggregatedCalculationLag) GetEx() (result AggregatedCalculationLagEx) {
	result.AggregatedCalculationLag = lag.Clone()
	result.Count = 1
	result.Sum = float64(lag.Average)
	return
}

func (lag *AggregatedCalculationLag) Clone() (result AggregatedCalculationLag) {
	result.Min = lag.Min
	result.Average = lag.Average
	result.Max = lag.Max
	return
}

func (lag *AggregatedCalculationLag) GetAllHours() []float64 {
	return []float64{lag.Min.Hours(), lag.Average.Hours(), lag.Max.Hours()}
}

type AggregatedCalculationLagEx struct {
	AggregatedCalculationLag
	Count int
	Sum   float64
}

func (lag *AggregatedCalculationLagEx) Aggregate(item AggregatedCalculationLagEx) {
	if item.Min < lag.Min {
		lag.Min = item.Min
	}
	if lag.Max < item.Max {
		lag.Max = item.Max
	}
	lag.Count += item.Count
	lag.Sum += item.Sum
}

func (lag *AggregatedCalculationLagEx) FinalizeAggregation() AggregatedCalculationLag {
	if lag.Count > 0 {
		lag.Average = time.Duration(
			int64(lag.Sum / float64(lag.Count)),
		)
	}
	return lag.AggregatedCalculationLag
}

func (lag *AggregatedCalculationLagEx) Clone() (result AggregatedCalculationLagEx) {
	result.AggregatedCalculationLag = lag.AggregatedCalculationLag.Clone()
	result.Count = lag.Count
	result.Sum = lag.Sum
	return
}
