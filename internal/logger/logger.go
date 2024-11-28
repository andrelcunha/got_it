package logger

import (
	"log"
)

type Logger struct {
	verbose bool
	debug   bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

func (l *Logger) Log(format string, args ...interface{}) {
	if l.verbose {
		log.Printf(format, args...)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		log.Printf(format, args...)
	}
}
