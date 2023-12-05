package pgxmigrations

import (
	"testing"
)

func TestSetQueriesFromDir(t *testing.T) {
	t.Run("invalid query files", func(t *testing.T) {
		migrator := NewMigrator("./testdata/nonexistent", nil)

		err := migrator.setQueriesFromDir()
		if err == nil {
			t.Fatalf("unexpected error got <nil>")
		}
	})

	t.Run("invalid query files", func(t *testing.T) {
		migrator := NewMigrator("./testdata/valid", nil)

		err := migrator.setQueriesFromDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		const expectedQueries = 1
		if len(migrator.inst.queries) != expectedQueries {
			t.Fatalf("expected %d queries, got %d", expectedQueries, len(migrator.inst.queries))
		}
	})

	t.Run("no files, just sub directories", func(t *testing.T) {
		migrator := NewMigrator("./testdata/subdirs", nil)
		err := migrator.setQueriesFromDir()
		if err == nil {
			t.Fatal("expected error for directory, got <nil>")
		}
	})

	t.Run("invalid file name", func(t *testing.T) {
		migrator := NewMigrator("./testdata/invalidfiles", nil)
		err := migrator.setQueriesFromDir()
		if err == nil {
			t.Fatal("expected error for incorrect file naming, got nil")
		}
	})
}
