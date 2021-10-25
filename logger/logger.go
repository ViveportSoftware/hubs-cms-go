package logger

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Info writes info logs
var Info *log.Logger

// Warn writes warn logs
var Warn *log.Logger

// Debug writes debug logs
var Debug *log.Logger

// Error writes error logs
var Error *log.Logger

// Setup setup logger functions
func Setup(logLevel string) {

	logLevel = strings.ToLower(logLevel)
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	debugHandle := ioutil.Discard
	errorHandle := os.Stderr

	if logLevel == "info" {
		infoHandle = os.Stdout
	}

	if logLevel == "info" || logLevel == "warn" {
		warnHandle = os.Stdout
	}

	if logLevel == "info" || logLevel == "warn" || logLevel == "debug" {
		debugHandle = os.Stdout
	}

	Info = log.New(infoHandle, "INF: ", log.Ldate|log.Lmicroseconds|log.LUTC)
	Warn = log.New(warnHandle, "WRN: ", log.Ldate|log.Lmicroseconds|log.LUTC)
	Debug = log.New(debugHandle, "DEB: ", log.Ldate|log.Lmicroseconds|log.LUTC)
	Error = log.New(errorHandle, "ERR: ", log.Ldate|log.Lmicroseconds|log.LUTC)
}
