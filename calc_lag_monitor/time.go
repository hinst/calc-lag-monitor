package main

import "time"

type TimeMeasurementUnit int8

const (
	TimeMeasurementUnitNone   TimeMeasurementUnit = 0
	TimeMeasurementUnitSecond TimeMeasurementUnit = 1
	TimeMeasurementUnitMinute TimeMeasurementUnit = 2
	TimeMeasurementUnitHour   TimeMeasurementUnit = 3
	TimeMeasurementUnitDay    TimeMeasurementUnit = 4
	TimeMeasurementUnitMonth  TimeMeasurementUnit = 5
)

func TruncateTime(t time.Time, unit TimeMeasurementUnit) time.Time {
	switch unit {
	case TimeMeasurementUnitSecond:
		return t.Truncate(time.Second)
	case TimeMeasurementUnitMinute:
		return t.Truncate(time.Minute)
	case TimeMeasurementUnitHour:
		return t.Truncate(time.Hour)
	case TimeMeasurementUnitDay:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case TimeMeasurementUnitMonth:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	default:
		return t
	}
}
