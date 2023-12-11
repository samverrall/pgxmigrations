package pgxmigrations

import (
	"fmt"
	"log/slog"
)

// logger wraps a slog logger, and a enable flog.
// by default logging is disabled, and can be enabled by using
// the WithDebugLogging opt func when initialising a Migrator.
type logger struct {
	enabled bool
	logger  *slog.Logger
}

func newLogger(enabled bool) *logger {
	return &logger{
		enabled: enabled,
		logger:  slog.Default(),
	}
}

func (l *logger) enable() {
	l.enabled = true
}

func (l *logger) Debug(msg string, args ...any) {
	if !l.enabled {
		return
	}
	l.logger.Debug(fmt.Sprintf("[pgxmigrations]: %s", msg), args...)
}

func (l *logger) Info(msg string, args ...any) {
	if !l.enabled {
		return
	}
	l.logger.Info(fmt.Sprintf("[pgxmigrations]: %s", msg), args...)
}

func (l *logger) Error(msg string, args ...any) {
	if !l.enabled {
		return
	}
	l.logger.Error(fmt.Sprintf("[pgxmigrations]: %s", msg), args...)
}
