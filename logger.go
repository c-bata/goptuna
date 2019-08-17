package goptuna

import (
	"fmt"
	"io"
	"log"
)

// Logger is the interface for logging messages.
// If you need to print more verbose logs, please use
// StudyOptionSetTrialNotifyChannel option.
type Logger interface {
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// LoggerLevel represents the level of logging.
type LoggerLevel int

const (
	LoggerLevelDiscard LoggerLevel = iota
	LoggerLevelDebug
	LoggerLevelWarn
	LoggerLevelInfo
	LoggerLevelError
)

func NewStdLogger(out io.Writer) *StdLogger {
	return &StdLogger{
		Logger:    log.New(out, "", log.LstdFlags),
		Level:     LoggerLevelDebug,
		WithColor: true,
	}
}

var _ Logger = &StdLogger{}

type StdLogger struct {
	Logger    *log.Logger
	Level     LoggerLevel
	WithColor bool
}

func (l *StdLogger) write(msg string, fields ...interface{}) {
	if l.Logger == nil {
		return
	}
	fields = append([]interface{}{msg}, fields...)
	l.Logger.Println(fields...)
}

func (l *StdLogger) Debug(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelDiscard {
		return
	}

	prefix := "[DEBUG] "
	if l.WithColor {
		prefix = "\033[1;31m" + prefix + "\033[0m "
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}

func (l *StdLogger) Warn(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelWarn {
		return
	}

	prefix := "[WARN] "
	if l.WithColor {
		prefix = "\033[1;33m" + prefix + "\033[0m "
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}

func (l *StdLogger) Info(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelInfo {
		return
	}

	prefix := "[INFO] "
	if l.WithColor {
		prefix = "\033[1;34m" + prefix + "\033[0m "
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}

func (l *StdLogger) Error(msg string, fields ...interface{}) {
	if l.Level > LoggerLevelError {
		return
	}

	prefix := "[ERROR] "
	if l.WithColor {
		prefix = "\033[1;36m" + prefix + "\033[0m "
	}
	l.write(fmt.Sprintf("%s%s:", prefix, msg), fields...)
}
