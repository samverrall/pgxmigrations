package pgxmigrations

type OptFunc func(m *Migrator)

// WithLogging can be passed to `NewMigrator` to optionally enable logging.
// Logging is disabled by default.
func WithLogging(enableLogging bool) OptFunc {
	return func(m *Migrator) {
		m.logging = enableLogging
	}
}

func WithInternalTable(withInternalTable bool) OptFunc {
	return func(m *Migrator) {
		m.internalTable = withInternalTable
	}
}
