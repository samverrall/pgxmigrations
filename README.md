[![Go Reference](https://pkg.go.dev/badge/github.com/samverrall/pgmigrations.svg)](https://pkg.go.dev/github.com/samverrall/pgmigrations)

# pgmigrations - PostgresSQL Migrations in Go with The PGX Driver and Toolkit

## Example Usuage 

```go 
func main () {
	ctx := context.Background()

	migrator := pgxmigrations.NewMigrator("/migrations/dir", pgx.DB, pgxmigrations.WithDebugLogging(true))
	if err := migrator.Migrate(ctx); err != nil {
		fmt.Fatal("failed to migrate postgres migrations", err)
	}
}

```

### Features 
