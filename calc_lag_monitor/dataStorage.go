package main

import (
	"bytes"
	"time"

	bolt "go.etcd.io/bbolt"
)

type DataStorage struct {
	db *bolt.DB
}

const DATA_STORAGE_FILE_PATH = "./data.db"
const CALCULATION_LAG_INFO_ROW_BUCKET_NAME = "CalculationLagInfoRow"
const PERMISSION_EVERYBODY_READ_WRITE = 0666
const OUTPUT_ROW_COUNT_LIMIT = 1000

var CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES = []byte(CALCULATION_LAG_INFO_ROW_BUCKET_NAME)

func (storage *DataStorage) Open() {
	db, err := bolt.Open(DATA_STORAGE_FILE_PATH,
		PERMISSION_EVERYBODY_READ_WRITE,
		&bolt.Options{Timeout: 30 * time.Second})
	AssertWrapped(err, "Unable to open file "+DATA_STORAGE_FILE_PATH)
	storage.db = db
}

func (storage *DataStorage) SaveCalculationLagInfoRow(row *CalculationLagInfoRow) {
	storage.db.Update(func(transaction *bolt.Tx) error {
		bucket, bucketError := transaction.CreateBucketIfNotExists(
			CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
		AssertWrapped(bucketError, "Unable to create bucket")
		var buffer bytes.Buffer
		row.Write(&buffer)
		bucket.Put(Int64ToBytes(row.Time.UnixMilli()), buffer.Bytes())
		return nil
	})
}

func (storage *DataStorage) ReadCalculationLagInfoRows(startUnixMillis int64, endUnixMillis int64) (
	rows []*CalculationLagInfoRow,
) {
	rowMap := make(map[int64]*CalculationLagInfoRow)
	storage.db.View(func(transaction *bolt.Tx) error {
		cursor := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES).Cursor()
		var key []byte
		var value []byte
		if startUnixMillis != 0 {
			key, value = cursor.Seek(Int64ToBytes(startUnixMillis))
		} else {
			key, value = cursor.First()
		}
		for key != nil {
			if value != nil {
				row := &CalculationLagInfoRow{}
				row.Read(bytes.NewBuffer(value))
				rowMap[BytesToInt64(key)] = row
			}
			key, value = cursor.Next()
			if endUnixMillis != 0 && !(BytesToInt64(key) < endUnixMillis) {
				break
			}
		}
		return nil
	})
	return
}

func (storage *DataStorage) Close() {
	error := storage.db.Close()
	AssertWrapped(error, "Unable to close database")
}

type CalculationLagInfoRowResponseBuilder struct {
	// Inputs
	StartUnixMillis int64
	EndUnixMillis   int64

	AggregationLevel TimeMeasurementUnit
	// Outputs
	Rows map[int64]*CalculationLagInfoRow
}

func (builder *CalculationLagInfoRowResponseBuilder) Build(transaction *bolt.Tx) error {
	builder.Rows = make(map[int64]*CalculationLagInfoRow)
	cursor := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES).Cursor()
	var key []byte
	var value []byte
	if builder.StartUnixMillis != 0 {
		key, value = cursor.Seek(Int64ToBytes(builder.StartUnixMillis))
	} else {
		key, value = cursor.First()
	}
	for key != nil {
		if value != nil {
			row := &CalculationLagInfoRow{}
			row.Read(bytes.NewBuffer(value))
			builder.Rows[BytesToInt64(key)] = row
		}
		key, value = cursor.Next()
		if builder.EndUnixMillis != 0 && !(BytesToInt64(key) < builder.EndUnixMillis) {
			break
		}
	}
	return nil
}

func (builder *CalculationLagInfoRowResponseBuilder) addRow(row *CalculationLagInfoRow) {
	rowTime := TruncateTime(row.Time, builder.AggregationLevel).UnixMilli()
	builder.Rows[rowTime] = row
	for len(builder.Rows) > OUTPUT_ROW_COUNT_LIMIT {
		builder.AggregationLevel = builder.AggregationLevel.GetNextOrFail()
	}
}
