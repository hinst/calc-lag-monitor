package main

import "strconv"

func ParseIntOr0(text string) (result int64) {
	if len(text) > 0 {
		var parseError error
		result, parseError = strconv.ParseInt(text, 10, 64)
		AssertWrapped(parseError, "Cannot parse text as integer '"+text+"'")
	}
	return
}
