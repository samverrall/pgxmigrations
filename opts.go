package pgxmigrations

type OptFunc func(m *Migrator)

func WithDebugLogging(debugLogging bool) OptFunc {
	return func(m *Migrator) {
		m.debugLogging = debugLogging
	}
}

func WithInternalTable(withInternalTable bool) OptFunc {
	return func(m *Migrator) {
		m.internalTable = withInternalTable
	}
}
