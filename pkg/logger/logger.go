package logger

import (
	"fmt"
	"log"
)

const (
	LevelInfo = iota
	LevelDebug
)

type Logger struct {
	currentLogLevel int
}

func New(logLevelStr string) *Logger {
	logLevel := LevelDebug
	if logLevelStr == "info" {
		logLevel = LevelInfo
	}

	log.SetFlags(log.Ldate | log.Lmicroseconds)

	return &Logger{
		currentLogLevel: logLevel,
	}
}

func (l *Logger) Info(v ...any) {
	log.Println("[INFO]:", fmt.Sprint(v...))
}

func (l *Logger) Debug(v ...any) {
	if l.currentLogLevel >= LevelDebug {
		log.Println("[DEBUG]:", fmt.Sprint(v...))
	}
}
