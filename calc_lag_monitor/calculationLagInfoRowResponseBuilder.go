package main

import (
	"bytes"
	"sort"

	bolt "go.etcd.io/bbolt"
)

type CalculationLagInfoRowResponseBuilder struct {
	// Inputs
	StartUnixMillis int64
	EndUnixMillis   int64

	AggregationLevel TimeMeasurementUnit
	// Outputs
	Rows map[int64]*CalculationLagInfoRow
}

type CalculationLagAggregatedRows struct {
	Rows             []*CalculationLagInfoRow
	AggregationLevel TimeMeasurementUnit
}

func (builder *CalculationLagInfoRowResponseBuilder) Build(transaction *bolt.Tx) error {
	builder.Rows = make(map[int64]*CalculationLagInfoRow)
	bucket := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
	if nil == bucket {
		return nil
	}
	cursor := bucket.Cursor()
	var key []byte
	var value []byte
	if builder.StartUnixMillis != 0 {
		key, value = cursor.Seek(Int64ToBytes(builder.StartUnixMillis))
	} else {
		key, value = cursor.First()
	}
	for key != nil {
		keyIsInRange := (builder.StartUnixMillis <= 0 || builder.StartUnixMillis <= BytesToInt64(key)) &&
			(builder.EndUnixMillis <= 0 || BytesToInt64(key) < builder.EndUnixMillis)
		if value != nil && keyIsInRange {
			row := &CalculationLagInfoRow{}
			row.Read(bytes.NewBuffer(value))
			builder.addRow(row)
		}
		key, value = cursor.Next()
	}
	return nil
}

func (builder *CalculationLagInfoRowResponseBuilder) addRow(row *CalculationLagInfoRow) {
	rowTime := TruncateTime(row.Time, builder.AggregationLevel).UnixMilli()
	builder.Rows[rowTime] = row
	for len(builder.Rows) > OUTPUT_ROW_COUNT_LIMIT {
		builder.AggregationLevel = builder.AggregationLevel.GetNextOrFail()
		builder.collapseRows()
	}
}

func (builder *CalculationLagInfoRowResponseBuilder) collapseRows() {
	multiRows := make(map[int64][]*CalculationLagInfoRow)
	for rowTime, row := range builder.Rows {
		rowTime = TruncateTime(row.Time, builder.AggregationLevel).UnixMilli()
		multiRows[rowTime] = append(multiRows[rowTime], row)
	}
	builder.Rows = make(map[int64]*CalculationLagInfoRow)
	for rowTime, rows := range multiRows {
		builder.Rows[rowTime] = AggregateCalculationLagInfoRows(rows)
	}
}

func (builder *CalculationLagInfoRowResponseBuilder) GetRowArray() []*CalculationLagInfoRow {
	array := make([]*CalculationLagInfoRow, 0, len(builder.Rows))
	for _, item := range builder.Rows {
		array = append(array, item)
	}
	sort.Slice(array, func(i int, j int) bool {
		return array[i].Time.Before(array[j].Time)
	})
	return array
}

func (builder *CalculationLagInfoRowResponseBuilder) GetResponse() (result CalculationLagAggregatedRows) {
	result.Rows = builder.GetRowArray()
	result.AggregationLevel = builder.AggregationLevel
	return
}
