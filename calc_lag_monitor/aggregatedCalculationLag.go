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
