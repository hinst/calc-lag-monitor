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

func AggregateCalculationLagInfoRows(rows []*CalculationLagInfoRow) *CalculationLagInfoRow {
	if len(rows) <= 0 {
		return nil
	}
	var aggregator CalculationLagInfoRowEx
	aggregator.InitializeAggregation(rows[0])
	for _, item := range rows {
		aggregator.Aggregate(item)
	}
	result := aggregator.FinalizeAggregation()
	return &result
}

type CalculationLagInfoRowEx struct {
	CalculationLagInfoRow
	Cheap     AggregatedCalculationLagEx
	Expensive AggregatedCalculationLagEx
}

func (row *CalculationLagInfoRowEx) InitializeAggregation(firstItem *CalculationLagInfoRow) {
	row.Time = firstItem.Time
	row.Cheap.InitializeAggregation(firstItem.Cheap)
	row.Expensive.InitializeAggregation(firstItem.Expensive)
}

func (row *CalculationLagInfoRowEx) Aggregate(item *CalculationLagInfoRow) {
	if item.Time.Before(row.Time) {
		row.Time = item.Time
	}
	row.Cheap.Aggregate(item.Cheap)
	row.Expensive.Aggregate(item.Expensive)
}

func (row *CalculationLagInfoRowEx) FinalizeAggregation() CalculationLagInfoRow {
	row.CalculationLagInfoRow.Cheap = row.Cheap.FinalizeAggregation()
	row.CalculationLagInfoRow.Expensive = row.Expensive.FinalizeAggregation()
	return row.CalculationLagInfoRow
}
