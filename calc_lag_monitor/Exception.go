package main

import "fmt"

type Exception struct {
	cause   error
	message string
}

func CreateException(message string, cause error) Exception {
	return Exception{message: message, cause: cause}
}

func (exception Exception) Error() string {
	var result = exception.message
	if exception.cause != nil {
		result += "\n" + exception.cause.Error()
	}
	return result
}

func (exception Exception) String() string {
	return exception.Error()
}

var _ error = Exception{}
var _ fmt.Stringer = Exception{}

func AssertWrapped(e error, message string) {
	if e != nil {
		if len(message) > 0 {
			panic(CreateException(message, e))
		} else {
			panic(e)
		}
	}
}
