package main

import (
	"bytes"
	"encoding/binary"

	bolt "go.etcd.io/bbolt"
)

type DataMigration interface {
	Migrate(t *bolt.Tx)
}

type BigEndianMigration struct {
}

var DataMigrations = []DataMigration{BigEndianMigration{}}

func (migration BigEndianMigration) Migrate(t *bolt.Tx) {
	calculationLagInfoRowBucket := t.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
	if calculationLagInfoRowBucket != nil {
		var rows []CalculationLagInfoRow

		DEFAULT_ENCODING = binary.LittleEndian
		cursor := calculationLagInfoRowBucket.Cursor()
		key, value := cursor.First()
		for key != nil {
			var row CalculationLagInfoRow
			row.Read(bytes.NewBuffer(value))
			rows = append(rows, row)
			key, value = cursor.Next()
		}

		DEFAULT_ENCODING = binary.BigEndian
		t.DeleteBucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
		var bucketError error
		calculationLagInfoRowBucket, bucketError = t.CreateBucketIfNotExists(
			CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
		AssertWrapped(bucketError, "Unable to create bucket")
		for _, row := range rows {
			var buffer bytes.Buffer
			row.Write(&buffer)
			calculationLagInfoRowBucket.Put(Int64ToBytes(row.Time.UnixMilli()), buffer.Bytes())
		}
	}
}
