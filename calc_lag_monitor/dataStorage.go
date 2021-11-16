package main

import (
	"bytes"
	"time"

	bolt "go.etcd.io/bbolt"
)

type DataStorage struct {
	Configuration *Configuration
	db            *bolt.DB
}

type DataStorageStatistics struct {
	CountOfCalculationLagRecords int
}

const CALCULATION_LAG_INFO_ROW_BUCKET_NAME = "CalculationLagInfoRow"
const PERMISSION_EVERYBODY_READ_WRITE = 0666
const OUTPUT_ROW_COUNT_LIMIT = 2000

var CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES = []byte(CALCULATION_LAG_INFO_ROW_BUCKET_NAME)

func (storage *DataStorage) Open() {
	db, err := bolt.Open(storage.Configuration.BoltDbFilePath,
		PERMISSION_EVERYBODY_READ_WRITE,
		&bolt.Options{Timeout: 30 * time.Second})
	AssertWrapped(err, "Unable to open file "+storage.Configuration.BoltDbFilePath)
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

func (storage *DataStorage) ReadCalculationLagInfoRows(startUnixMillis int64, endUnixMillis int64,
) CalculationLagAggregatedRows {
	builder := &CalculationLagInfoRowsBuilder{
		StartUnixMillis: startUnixMillis, EndUnixMillis: endUnixMillis}
	storage.db.View(builder.Build)
	return builder.GetResponse()
}

func (storage *DataStorage) GetStatistics() (result DataStorageStatistics) {
	storage.db.View(func(transaction *bolt.Tx) error {
		cursor := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES).Cursor()
		key, _ := cursor.First()
		for key != nil {
			result.CountOfCalculationLagRecords += 1
			key, _ = cursor.Next()
		}
		return nil
	})
	return
}

func (storage *DataStorage) Close() {
	error := storage.db.Close()
	AssertWrapped(error, "Unable to close database")
}
