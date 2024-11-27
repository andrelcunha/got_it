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

func (l *Logger) Log(message string) {
	if l.verbose {
		log.Println(message)
	}
}
