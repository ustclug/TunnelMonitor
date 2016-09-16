package main

import (
	"github.com/alecthomas/log4go"
)

var logger log4go.Logger

func initLogger() {
	logger = make(log4go.Logger)
	logger.AddFilter("stdout", log4go.INFO, log4go.NewConsoleLogWriter())

	if logFileName := configCommon("log", INFO); logFileName != "" {
		logger.AddFilter("file", log4go.INFO, log4go.NewFileLogWriter(logFileName, false).SetFormat("[%D %T][%L]%M - %S"))
	}
	logger.Info("log module initialized")
}
