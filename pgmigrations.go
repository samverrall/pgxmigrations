package pgxmigrations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Migrator struct {
	dir                string
	logging            bool
	disableForeignKeys bool
	logger             *logger
	inst               *migratorInst
}

type migratorInst struct {
	queries []string
	db      DB
}

// NewMigrator returns a new migrator
func NewMigrator(migrationsDir string, db DB, opts ...OptFunc) *Migrator {
	m := &Migrator{
		dir:    migrationsDir,
		logger: newLogger(false),
		inst: &migratorInst{
			db:      db,
			queries: make([]string, 0),
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.logging {
		m.logger.enable()
		m.logger.Debug("debug logging enabled")
	}

	return m
}

func (m *Migrator) setQueriesFromDir() error {
	m.logger.Info("reading migrations from dir", "dir", m.dir)

	migrations, err := os.ReadDir(m.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("supplied migrations directory does not exist: %s", m.dir)
		}
		return fmt.Errorf("read migrations dir: %w", err)
	}

	queries := make([]string, len(migrations))
	for i, f := range migrations {
		filename := f.Name()

		if f.IsDir() {
			return fmt.Errorf("want file; got directory %q", filename)
		}

		if filename[:4] != fmt.Sprintf("%04d", i+1) {
			return fmt.Errorf("want file beginning with %04d; got %q", i+1, filename)
		}

		b, err := os.ReadFile(filepath.Join(m.dir, filename))
		if err != nil {
			return fmt.Errorf("read migration file: %w", err)
		}

		queries[i] = string(b)
	}
	m.inst.queries = queries
	return nil
}

// Migrate attempts to exec migration files found in your migrations
// directory.
func (m *Migrator) Migrate(ctx context.Context) error {
	err := m.setQueriesFromDir()
	if err != nil {
		return err
	}
	return m.migrate(ctx)
}

func (m *Migrator) migrate(ctx context.Context) error {
	migrations := m.inst.queries
	if len(migrations) == 0 {
		return nil
	}

	tx, err := m.inst.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			m.logger.Error("rollback", "error", err)
		}
	}()

	if err := createMigrationsTbl(ctx, tx); err != nil {
		return err
	}

	m.logger.Info("insert default 0 value into _migartions")
	_, err = tx.Exec(ctx, "INSERT INTO _migrations (version) VALUES (0) ON CONFLICT DO NOTHING")
	if err != nil {
		m.logger.Error("failed to insert intial migration value", "error", err)
		return fmt.Errorf("inserting initial 0 version: exec: %w", err)
	}

	var count int
	if err := tx.QueryRow(ctx, "SELECT COUNT(1) FROM _migrations;").Scan(&count); err != nil {
		m.logger.Error("select count(1) from _migrations error", "error", err)
		return fmt.Errorf("failed to query for count: exec: %w", err)
	}

	var version int
	stmt := "SELECT version FROM _migrations;"
	if err := tx.QueryRow(ctx, stmt).Scan(&version); err != nil {
		m.logger.Error("select version from _migrations error", "error", err)
		return err
	}

	if version < 0 {
		return fmt.Errorf("want current migration to be version 0 or more; got %v", version)
	}
	migration := version
	nm := len(migrations)

	// If the number of migration strings is less than the version then we must have
	// lost some migrations and the data cannot be trusted
	if nm < version {
		return fmt.Errorf("want at least %v migration strings; got %v", version, nm)
	}

	// If the version is the same as the number of migration strings then we must be up to date
	if nm == version {
		return nil
	}

	if m.disableForeignKeys {
		if err := disableForeignKeys(ctx, tx); err != nil {
			m.logger.Error("failed to disable foreign keys", "error", err)
			return fmt.Errorf("failed to disable foreign keys: %w", err)
		}
	}

	for i, stmt := range migrations {
		if i < version {
			continue
		}

		// If the migration file is empty then don't waste the
		// time trying to execute a query
		if stmt = strings.TrimSpace(stmt); stmt == "" {
			continue
		}
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("migration: %d: %w", i, err)
		}

		migration++
	}

	if m.disableForeignKeys {
		if err := enableForeignKeys(ctx, tx); err != nil {
			m.logger.Error("failed to enable foreign keys", "error", err)
			return fmt.Errorf("failed to enable foreign keys: %w", err)
		}
	}

	// If the migration number is greater than the starting version then
	// that means we must have executed some migration strings so we
	// should attempt to set the migration version to the new number
	if migration > version {
		stmt := `
			UPDATE _migrations SET version = $1, updated_at = $2 
		`
		args := []any{
			migration,
			time.Now().UTC(),
		}
		if _, err := tx.Exec(ctx, stmt, args...); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		m.logger.Error("commit", "error", err)
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func disableForeignKeys(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, "SET session_replication_role = replica;")
	if err != nil {
		return err
	}
	return nil
}

func enableForeignKeys(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, "SET session_replication_role = origin;")
	if err != nil {
		return err
	}
	return nil
}

func createMigrationsTbl(ctx context.Context, tx pgx.Tx) error {
	stmt := `
CREATE TABLE IF NOT EXISTS _migrations (
    version    INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
	`

	if _, err := tx.Exec(ctx, stmt); err != nil {
		return err
	}

	return nil
}
