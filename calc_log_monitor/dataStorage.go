package main

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

type DataStorage struct {
	db *bolt.DB
}

const filePath = "./data.db"

func (data *DataStorage) Open() {
	db, err := bolt.Open(filePath, 0666, &bolt.Options{Timeout: 30 * time.Second})
	AssertWrapped(err, "Unable to open file "+filePath)
	data.db = db
}

func (data *DataStorage) SaveCalculationLagInfoRow() {
}
