package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *SimpleLogger) Info(msg string, keysAndValues ...interface{}) {
	l.infoLogger.Println(msg, keysAndValues)
}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	l.errorLogger.Println(msg, keysAndValues)
}
