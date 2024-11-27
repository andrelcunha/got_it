package logger

import (
	"log"
)

type Logger struct {
	verbose bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

func (l *Logger) Log(format string, args ...interface{}) {
	if l.verbose {
		log.Printf(format, args...)
	}
}
