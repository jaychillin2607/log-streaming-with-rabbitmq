package main

import (
	"log"
	"os"
)

// enum
type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	WARNING
	ERROR
)

type logMessage struct {
	Level LogLevel `json:"level"`
	Msg   string   `json:"msg"`
}

func NewLogMessage(level, msg *string) *logMessage {
	message := new(logMessage)
	message.Msg = *msg
	message.Level = castStrToLogLevel(level)
	return message
}

type multiLogger struct {
	info    *log.Logger
	debug   *log.Logger
	warning *log.Logger
	error   *log.Logger
}

func (ml *multiLogger) Info(v string)    { ml.info.Println(v) }
func (ml *multiLogger) Debug(v string)   { ml.debug.Println(v) }
func (ml *multiLogger) Warning(v string) { ml.warning.Println(v) }
func (ml *multiLogger) Error(v string)   { ml.error.Println(v) }

func newMultiLogger() *multiLogger {
	return &multiLogger{
		info:    log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		debug:   log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		warning: log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile),
		error:   log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func getEnvOrDefault(envVar, defaultVar string) string {
	if val := os.Getenv(envVar); val != "" {
		return val
	}
	return defaultVar
}

func handleFailure(errMsg string, err *error) {
	if *err != nil {
		log.Fatalf("Message: %s || Error: %v", errMsg, *err)
	}
}

func castStrToLogLevel(level *string) LogLevel {
	switch *level {
	case "INFO":
		return INFO
	case "DEBUG":
		return DEBUG
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	default:
		handleFailure("Invalid log level provided. Valid levels [INFO, DEBUG, WARNING, ERROR]", &EMPTY_ERROR)
		return -1
	}
}

func createExchangeKey(level *string) string {
	switch *level {
	case "INFO":
		return "INFO.*.*"
	case "DEBUG":
		return "*.DEBUG.*"
	case "WARNING":
		return "*.*.WARNING"
	case "ERROR":
		return "ERROR.*"
	default:
		handleFailure("Invalid log level provided. Valid levels [INFO, DEBUG, WARNING, ERROR]", &EMPTY_ERROR)
		return DEFAULT_TOPIC
	}
}
