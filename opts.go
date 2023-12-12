package pgxmigrations

type OptFunc func(m *Migrator)

// WithLogging can be passed to `NewMigrator` to optionally enable logging.
// Logging is disabled by default.
func WithLogging(enableLogging bool) OptFunc {
	return func(m *Migrator) {
		m.logging = enableLogging
	}
}

// WithDisableForeignKeys configures the Migrator to disable
// foreign key constraints during the migration process if the provided 'disable'
// parameter is set to true.
func WithDisableForeignKeys(disable bool) OptFunc {
	return func(m *Migrator) {
		m.disableForeignKeys = disable
	}
}
