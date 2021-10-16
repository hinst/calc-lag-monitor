package main

import "time"

type AggregatedCalculationRequest struct {
	Min     time.Time
	Average time.Time
	Max     time.Time
	Count   int
}
