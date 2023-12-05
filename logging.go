package pgxmigrations

import "log/slog"

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
	}
}

func (l *logger) enable() {
	l.enabled = true
}

func (l *logger) Debug(msg string, args ...any) {
	if !l.enabled {
		return
	}
	l.logger.Debug(msg, args...)
}
