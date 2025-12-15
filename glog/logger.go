package glog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type _LogLevel struct {
	Code  int
	Label string
}

var DEBUG = "DEBUG"
var INFO = "INFO"
var WARN = "WARN"
var ERROR = "ERROR"

var LogLevelMap = map[string]int{DEBUG: 1,
	INFO:  2,
	WARN:  3,
	ERROR: 4,
}

var LogLevel = _LogLevel{}

func (l *_LogLevel) Load() {
	l.Code = LogLevelMap[l.Label]
}

var _log *log.Logger = log.New(os.Stdout, "", 0)

func Log(v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_log.Println(timestamp, v)
}

func concatenateWithSpace(values ...interface{}) string {
	// Convert variadic values to string slice
	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = fmt.Sprint(v)
	}

	// Join string values with space
	result := strings.Join(stringValues, " ")
	return result
}

func LogL(level string, v ...interface{}) {
	if LogLevelMap[level] < LogLevel.Code {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_log.Println(timestamp, concatenateWithSpace(v...))
}

// LogFields logs a structured message with key=value pairs built from fields.
// Example:
//   glog.LogFields(glog.INFO, "http request",
//     map[string]interface{}{"method": "GET", "url": "/api", "status": 200})
func LogFields(level string, msg string, fields map[string]interface{}) {
	if LogLevelMap[level] < LogLevel.Code {
		return
	}

	parts := make([]string, 0, len(fields)+1)
	parts = append(parts, msg)
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_log.Println(timestamp, strings.Join(parts, " "))
}
