package mylog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Logger struct {
	mu     sync.Mutex
	file   *os.File
	logger *log.Logger
}

func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return &Logger{
		file:   f,
		logger: log.New(f, "", log.LstdFlags),
	}, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, fmt.Sprintf("[INFO] "+format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, fmt.Sprintf("[WARN] "+format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
}

func (l *Logger) Close() error {
	return l.file.Close()
}
