package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type BinaryObjectVersionNumber = uint8

var DEFAULT_ENCODING = binary.LittleEndian

func BinaryWrite(buffer *bytes.Buffer, data interface{}) {
	error := binary.Write(buffer, DEFAULT_ENCODING, data)
	AssertWrapped(error, "Unable to write binary data")
}

func BinaryRead(buffer *bytes.Buffer, data interface{}) {
	error := binary.Read(buffer, DEFAULT_ENCODING, data)
	AssertWrapped(error, "Unable to read binary data")
}

func Int64ToBytes(value int64) (result []byte) {
	result = make([]byte, 8)
	DEFAULT_ENCODING.PutUint64(result, uint64(value))
	return
}

func BytesToInt64(value []byte) int64 {
	return int64(DEFAULT_ENCODING.Uint64(value))
}

func EncodeJson(data interface{}) []byte {
	text, error := json.Marshal(data)
	AssertWrapped(error, "Unable to encode data as json")
	return text
}
