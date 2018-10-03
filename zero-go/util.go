package main

import "time"

// NullTime represent a sqlite time which is nullable
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// If is a util function for ?:
func If(condition bool, trueRet interface{}, falseRet interface{}) interface{} {
	if condition {
		return trueRet
	}
	return falseRet
}
