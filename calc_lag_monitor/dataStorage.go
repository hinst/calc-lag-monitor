package main

import (
	"bytes"
	"log"
	"strconv"
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
const MIGRATION_BUCKET_NAME = "Migrations"
const PERMISSION_EVERYBODY_READ_WRITE = 0666
const OUTPUT_ROW_COUNT_LIMIT = 2000

var CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES = []byte(CALCULATION_LAG_INFO_ROW_BUCKET_NAME)
var MIGRATION_BUCKET_NAME_BYTES = []byte(MIGRATION_BUCKET_NAME)

func (storage *DataStorage) Open() {
	db, err := bolt.Open(storage.Configuration.BoltDbFilePath,
		PERMISSION_EVERYBODY_READ_WRITE,
		&bolt.Options{Timeout: 30 * time.Second})
	AssertWrapped(err, "Unable to open file "+storage.Configuration.BoltDbFilePath)
	storage.db = db
	storage.Migrate()
}

func (storage *DataStorage) Migrate() {
	storage.db.Update(func(t *bolt.Tx) error {
		t.CreateBucketIfNotExists(MIGRATION_BUCKET_NAME_BYTES)
		migrations := t.Bucket(MIGRATION_BUCKET_NAME_BYTES)
		for index, migration := range DataMigrations {
			migrationKey := Int64ToBytes(int64(index))
			if migrations.Get(migrationKey) == nil {
				log.Println("Migrating " + strconv.Itoa(index))
				migration.Migrate(t)
				AssertWrapped(
					migrations.Put(migrationKey, []byte{1}),
					"Unable to store migration record")
			}
		}
		return nil
	})
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
		bucket := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
		if bucket != nil {
			cursor := bucket.Cursor()
			key, _ := cursor.First()
			for key != nil {
				result.CountOfCalculationLagRecords += 1
				key, _ = cursor.Next()
			}
		}
		return nil
	})
	return
}

func (storage *DataStorage) RemoveAnomalies(writeEnabled bool) (result string) {
	anomalyKeys := make([][]byte, 0)
	storage.db.View(func(transaction *bolt.Tx) error {
		bucket := transaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)

		if bucket != nil {
			cursor := bucket.Cursor()
			key, value := cursor.First()
			for key != nil {
				keyTime := time.UnixMilli(BytesToInt64(key))
				var row CalculationLagInfoRow
				row.Read(bytes.NewBuffer(value))
				rowTime := row.Time
				if !keyTime.Equal(rowTime) {
					result += "Inconsistency: key time does not match row time " +
						keyTime.String() + " " + row.Time.String() + "\n"
				}
				allHours := make([]float64, 0)
				allHours = append(allHours, row.Cheap.GetAllHours()...)
				allHours = append(allHours, row.Expensive.GetAllHours()...)
				var isAnomaly bool
				for _, hours := range allHours {
					if hours > 100_000 {
						isAnomaly = true
					}
				}
				if isAnomaly {
					result += "Inconsistency: duration is too long at " + row.Time.String() + "\n"
					anomalyKeys = append(anomalyKeys, key)
				}

				key, value = cursor.Next() // must be at the end of the loop
			}
		}
		return nil
	})
	result += "Total count of anomalies: " + strconv.Itoa(len(anomalyKeys)) + "\n"
	if writeEnabled {
		updateResult := storage.db.Update(func(t *bolt.Tx) error {
			calculationLagInfoBucket := t.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
			if calculationLagInfoBucket != nil {
				for _, anomalyKey := range anomalyKeys {
					calculationLagInfoBucket.Delete(anomalyKey)
				}
			}
			return nil
		})
		const errorMessage = "Update failed. Unable to remove anomalies"
		if updateResult != nil {
			result += errorMessage + "\n" + updateResult.Error()
		}
		AssertWrapped(updateResult, errorMessage)
		result += "Anomalies were removed: " + strconv.Itoa(len(anomalyKeys))
	}
	return
}

func (storage *DataStorage) Close() {
	error := storage.db.Close()
	AssertWrapped(error, "Unable to close database")
}
