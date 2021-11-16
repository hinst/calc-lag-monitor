package main

import (
	"strconv"
	"strings"
)

func ParseIntOr0(text string) (result int64) {
	if len(text) > 0 {
		var parseError error
		result, parseError = strconv.ParseInt(text, 10, 64)
		AssertWrapped(parseError, "Cannot parse text as integer '"+text+"'")
	}
	return
}

func BytesToString(bytes []byte) string {
	textParts := make([]string, len(bytes))
	for index, value := range bytes {
		textParts[index] = strconv.Itoa(int(value))
	}
	return strings.Join(textParts, " ")
}
