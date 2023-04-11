// Package logger implements wraps https://pkg.go.dev/log to instantiate a Syringe-compatible Logger.
package logger

import (
	"log"
	"os"
)

// Logger is the receiver.
type Logger struct {
	handler       *log.Logger // underlying standard log instance
	loggingToFile bool        // true when appending to a file
	file          *os.File    // file to append
}

// New instantiates a Logger. The destination may be "stderr" or "" (output to stderr), "stdout" (output to stdout), or a file to append.
func New(dst string) (*Logger, error) {
	l := &Logger{
		handler: log.Default(),
	}
	switch dst {
	case "stderr":
	case "":
		l.handler.SetOutput(os.Stderr)
	case "stdout":
		l.handler.SetOutput(os.Stdout)
	default:
		var err error
		l.file, err = os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		l.handler.SetOutput(l.file)
	}
	return l, nil
}

// Close shuts down the logger.
func (l *Logger) Close() {
	if l.loggingToFile {
		l.file.Close()
	}
}

// Print satisfies the syringe.Logger interface.
func (l *Logger) Print(v ...interface{}) {
	l.handler.Print(v...)
}
