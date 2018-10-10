package util

import (
	"log"
	"runtime"
	"time"
)

/// StrType is string type alias for `context.Context`
type StrType string

const (
	// LoadersKey for loaders
	LoadersKey StrType = "loaders"

	// ClientKey for client object
	ClientKey StrType = "client"
)

// HandleError prints error's line number
func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", fn, line, err)
		b = true
	}
	return
}

// ULog for log messages and other code related stuff
func ULog(msg interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[info] %s:%d %+v", fn, line, msg)
}

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
