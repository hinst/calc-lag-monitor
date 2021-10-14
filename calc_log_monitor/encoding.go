package main

import (
	"bytes"
	"encoding/binary"
)

type BinaryObjectVersionNumber = uint16

func BinaryWrite(buffer *bytes.Buffer, data interface{}) {
	AssertWrapped(binary.Write(buffer, binary.LittleEndian, data), "Unable to write binary data")
}
