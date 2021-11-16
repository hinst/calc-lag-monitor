package main

import (
	"bytes"
	"sort"

	bolt "go.etcd.io/bbolt"
)

type CalculationLagInfoRowsBuilder struct {
	// Inputs
	StartUnixMillis int64
	EndUnixMillis   int64

	AggregationLevel TimeMeasurementUnit
	// Outputs
	Rows map[int64]*CalculationLagInfoRowEx
}

type CalculationLagAggregatedRows struct {
	Rows             []*CalculationLagInfoRow
	AggregationLevel TimeMeasurementUnit
}

func (builder *CalculationLagInfoRowsBuilder) Build(transaction *bolt.Tx) error {
	builder.Rows = make(map[int64]*CalculationLagInfoRowEx)
	bucket := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
	if nil == bucket {
		return nil
	}
	cursor := bucket.Cursor()
	var key []byte
	var value []byte
	key, value = cursor.First()
	for key != nil && value != nil {
		keyIsInRange := (builder.StartUnixMillis <= 0 || builder.StartUnixMillis <= BytesToInt64(key)) &&
			(builder.EndUnixMillis <= 0 || BytesToInt64(key) < builder.EndUnixMillis)
		if keyIsInRange {
			row := &CalculationLagInfoRow{}
			row.Read(bytes.NewBuffer(value))
			builder.addRow(row)
		}
		key, value = cursor.Next()
	}
	return nil
}

func (builder *CalculationLagInfoRowsBuilder) addRow(row *CalculationLagInfoRow) {
	rowTime := TruncateTime(row.Time, builder.AggregationLevel).UnixMilli()
	existingRow := builder.Rows[rowTime]
	if existingRow != nil {
		existingRow.Aggregate(row.GetExPtr())
	} else {
		builder.Rows[rowTime] = row.GetExPtr()
	}
	for len(builder.Rows) > OUTPUT_ROW_COUNT_LIMIT {
		builder.AggregationLevel = builder.AggregationLevel.GetNextOrFail()
		builder.collapseRows()
	}
}

func (builder *CalculationLagInfoRowsBuilder) collapseRows() {
	multiRows := make(map[int64][]*CalculationLagInfoRowEx)
	for rowTime, row := range builder.Rows {
		rowTime = TruncateTime(row.Time, builder.AggregationLevel).UnixMilli()
		multiRows[rowTime] = append(multiRows[rowTime], row)
	}
	builder.Rows = make(map[int64]*CalculationLagInfoRowEx)
	for rowTime, rows := range multiRows {
		builder.Rows[rowTime] = AggregateCalculationLagInfoRows(rows)
	}
}

func (builder *CalculationLagInfoRowsBuilder) GetRowArray() []*CalculationLagInfoRow {
	array := make([]*CalculationLagInfoRow, 0, len(builder.Rows))
	for _, rowEx := range builder.Rows {
		array = append(array, rowEx.FinalizeAggregation())
	}
	sort.Slice(array, func(i int, j int) bool {
		return array[i].Time.Before(array[j].Time)
	})
	return array
}

func (builder *CalculationLagInfoRowsBuilder) GetResponse() (result CalculationLagAggregatedRows) {
	result.Rows = builder.GetRowArray()
	result.AggregationLevel = builder.AggregationLevel
	return
}
