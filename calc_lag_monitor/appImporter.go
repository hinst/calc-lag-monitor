package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

type AppImporter struct {
	DestinationDbFilePath string
	SourceDbFilePath      string
}

func (importer *AppImporter) Run() {
	if len(importer.SourceDbFilePath) == 0 {
		panic(errors.New("need source db file path"))
	}
	mainDb, mainDbError := bolt.Open(importer.DestinationDbFilePath,
		PERMISSION_EVERYBODY_READ_WRITE,
		&bolt.Options{Timeout: 30 * time.Second})
	AssertWrapped(mainDbError, "Unable to open file "+importer.DestinationDbFilePath)

	sourceDb, sourceDbError := bolt.Open(importer.SourceDbFilePath,
		PERMISSION_EVERYBODY_READ_WRITE,
		&bolt.Options{Timeout: 30 * time.Second})
	AssertWrapped(sourceDbError, "Unable to open file "+importer.SourceDbFilePath)

	mainDb.Update(func(mainTransaction *bolt.Tx) error {
		mainBucket, bucketError := mainTransaction.CreateBucketIfNotExists(
			CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
		AssertWrapped(bucketError, "Unable to create bucket")
		sourceDb.View(func(sourceTransaction *bolt.Tx) error {
			sourceBucket := sourceTransaction.Bucket(CALCULATION_LAG_INFO_ROW_BUCKET_NAME_BYTES)
			if sourceBucket != nil {
				cursor := sourceBucket.Cursor()
				key, value := cursor.First()
				counter := 0
				for key != nil {
					if value != nil {
						mainBucket.Put(key, value)
						counter++
					}
					key, value = cursor.Next()
				}
				log.Println("Imported records: " + strconv.Itoa(counter))
			} else {
				log.Println("Source database lacks calculation lag info bucket")
			}
			return nil
		})
		return nil
	})
}
