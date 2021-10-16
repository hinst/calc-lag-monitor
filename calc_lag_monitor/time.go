package main

import (
	"errors"
	"strconv"
	"time"
)

type TimeMeasurementUnit int8

const (
	TimeMeasurementUnitNone   TimeMeasurementUnit = 0
	TimeMeasurementUnitSecond TimeMeasurementUnit = 1
	TimeMeasurementUnitMinute TimeMeasurementUnit = 2
	TimeMeasurementUnitHour   TimeMeasurementUnit = 3
	TimeMeasurementUnitDay    TimeMeasurementUnit = 4
	TimeMeasurementUnitMonth  TimeMeasurementUnit = 5
	TimeMeasurementUnitYear   TimeMeasurementUnit = 6
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

func (unit TimeMeasurementUnit) GetNext() TimeMeasurementUnit {
	return unit + 1
}

func (unit TimeMeasurementUnit) GetNextOrFail() TimeMeasurementUnit {
	unit = unit.GetNext()
	if !unit.IsValid() {
		panic(errors.New("Unable to get next time measurement unit " + strconv.Itoa(int(unit))))
	}
	return unit
}

func (unit TimeMeasurementUnit) IsValid() bool {
	return TimeMeasurementUnitNone <= unit && unit <= TimeMeasurementUnitYear
}
