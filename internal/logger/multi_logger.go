package logger

import (
	"log"
	"os"

	"github.com/jaychillin2607/log-streaming-with-rabbitmq/internal"
	"github.com/jaychillin2607/log-streaming-with-rabbitmq/pkg/types"
)

// enum

func NewLogMessage(level, msg *string) *types.LogMessage {
	message := new(types.LogMessage)
	message.Msg = *msg
	message.Level = internal.CastStrToLogLevel(level)
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

func NewMultiLogger() *multiLogger {
	return &multiLogger{
		info:    log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		debug:   log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		warning: log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile),
		error:   log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
