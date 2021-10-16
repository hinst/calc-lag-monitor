package main

import (
	"bytes"
	"time"
)

type CalculationLagInfoRow struct {
	Time      time.Time
	Cheap     AggregatedCalculationLag
	Expensive AggregatedCalculationLag
}

func (row *CalculationLagInfoRow) Write(buffer *bytes.Buffer) {
	BinaryWrite(buffer, BinaryObjectVersionNumber(1))
	BinaryWrite(buffer, row.Time.UnixMilli())
	row.Cheap.Write(buffer)
	row.Expensive.Write(buffer)
}
