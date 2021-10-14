package main

import (
	"bytes"
	"time"
)

type CalculationLagInfoRow struct {
	Moment    time.Time
	Cheap     AggregatedCalculationLag
	Expensive AggregatedCalculationLag
}

func (row *CalculationLagInfoRow) Write(buffer *bytes.Buffer) {
	BinaryWrite(buffer, BinaryObjectVersionNumber(1))
	BinaryWrite(buffer, row.Moment.UnixMilli())
	row.Cheap.Write(buffer)
	row.Expensive.Write(buffer)
}
