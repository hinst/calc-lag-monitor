package main

import "time"

type CalculationRequest struct {
	CalculateAt time.Time `bson:"calculateAt"`
	DisplayName string    `bson:"displayName"`
}
