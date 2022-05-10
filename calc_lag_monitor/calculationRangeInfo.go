package main

// Measured in milliseconds
type CalculationRangeInfoItem struct {
	Min     int64
	Average int64
	Median  int64
	Max     int64
}

type CalculationRangeInfo struct {
	Start    CalculationRangeInfoItem
	End      CalculationRangeInfoItem
	Duration CalculationRangeInfoItem
	Count    int64
}
