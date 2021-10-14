package main

import (
	"bytes"
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

func (lag *AggregatedCalculationLag) ReadFromRequest(request *AggregatedCalculationRequest) {
	now := time.Now()
	lag.Min = now.Sub(request.Min)
	lag.Average = now.Sub(request.Average)
	lag.Max = now.Sub(request.Max)
}
