package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
	file   *os.File
}

func NewLogger() *Logger {
	f, err := os.OpenFile("brutestick.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	return &Logger{
		logger: log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile),
		file:   f,
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.SetPrefix("INFO: ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logger.SetPrefix("WARN: ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.SetPrefix("ERROR: ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Close() {
	l.file.Close()
}
