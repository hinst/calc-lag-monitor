package main

import (
	"bytes"
	"errors"
	"strconv"
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

func (row *CalculationLagInfoRow) Read(buffer *bytes.Buffer) {
	var version BinaryObjectVersionNumber
	BinaryRead(buffer, &version)
	if version != 1 {
		panic(errors.New("Expected version 1 but got " + strconv.Itoa(int(version))))
	}
	var timeUnixMilli int64
	BinaryRead(buffer, &timeUnixMilli)
	row.Time = time.UnixMilli(timeUnixMilli)
	row.Cheap.Read(buffer)
	row.Expensive.Read(buffer)
}
