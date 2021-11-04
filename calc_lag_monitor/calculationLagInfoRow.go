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

func (row *CalculationLagInfoRow) String() string {
	return row.Time.Format(time.RFC3339Nano) + " " +
		row.Cheap.Average.String() + " " + row.Expensive.Average.String()
}

func AggregateCalculationLagInfoRows(rows []*CalculationLagInfoRowEx) *CalculationLagInfoRowEx {
	if len(rows) <= 0 {
		return nil
	}
	var aggregator *CalculationLagInfoRowEx
	for _, item := range rows {
		if aggregator == nil {
			aggregator = item.ClonePtr()
		} else {
			aggregator.Aggregate(item)
		}
	}
	return aggregator
}

func (row *CalculationLagInfoRow) GetEx() (result CalculationLagInfoRowEx) {
	result.CalculationLagInfoRow = row.Clone()
	result.Cheap.AggregatedCalculationLag = row.Cheap.Clone()
	result.Expensive.AggregatedCalculationLag = row.Expensive.Clone()
	return
}

func (row *CalculationLagInfoRow) GetExPtr() *CalculationLagInfoRowEx {
	result := row.GetEx()
	return &result
}

func (row *CalculationLagInfoRow) Clone() (result CalculationLagInfoRow) {
	result.Time = row.Time
	result.Cheap = row.Cheap.Clone()
	result.Expensive = row.Expensive.Clone()
	return
}

func (row *CalculationLagInfoRow) ClonePtr() *CalculationLagInfoRow {
	result := row.Clone()
	return &result
}

type CalculationLagInfoRowEx struct {
	CalculationLagInfoRow
	Cheap     AggregatedCalculationLagEx
	Expensive AggregatedCalculationLagEx
}

func (row *CalculationLagInfoRowEx) Aggregate(item *CalculationLagInfoRowEx) {
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

func (row *CalculationLagInfoRowEx) Clone() (result CalculationLagInfoRowEx) {
	result.CalculationLagInfoRow = row.CalculationLagInfoRow.Clone()
	result.Cheap = row.Cheap.Clone()
	result.Expensive = row.Expensive.Clone()
	return
}

func (row *CalculationLagInfoRowEx) ClonePtr() *CalculationLagInfoRowEx {
	result := row.Clone()
	return &result
}
