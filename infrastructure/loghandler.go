package infrastructure

import (
	"io"
	"log"
)

type logger interface {
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// NewLog - Create logger for project
func NewLog(out io.Writer, prefix string, flag int) *logHandler {
	l := logHandler{}
	l.logger = log.New(out, prefix, flag)
	return &l
}

type logHandler struct {
	logger *log.Logger
}

// Fatalf - Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (l *logHandler) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

func (l *logHandler) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}
