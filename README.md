[![Go Reference](https://pkg.go.dev/badge/github.com/samverrall/pgxmigrations.svg)](https://pkg.go.dev/github.com/samverrall/pgxmigrations)
[![GitHub License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/samverrall/pgxmigrations/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/release/samverrall/pgxmigrations.svg)](https://github.com/samverrall/pgxmigrations/releases)

# pgmigrations - PostgresSQL Migrations in Go with The PGX Driver and Toolkit

## Example Usage 

```go 
func main () {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, postgresConnectionURL)
	if err != nil {
		log.Fatal("failed to connect to postgres", err)
	}

	migrator := pgxmigrations.NewMigrator("/path/to/migrations", conn, pgxmigrations.WithLogging(true))
	if err := migrator.Migrate(ctx); err != nil {
		log.Fatal("failed to migrate postgres migrations", err)
	}
}
```

## Install

```
go get -u github.com/samverrall/pgxmigrations@latest
```

### ⚠️ Pre-release Notice ⚠️

**Important: This package is currently in pre-release and is not recommended for production use.**

#### What does "pre-release" mean?

This software is still in the early stages of development. While I strive to provide a stable and functional experience, there may be bugs, incomplete features, or breaking changes in future releases.

#### Why use a pre-release version?

- **Early Testing:** Get a glimpse of the upcoming features and improvements.
- **Community Feedback:** I encourage users to provide feedback, report issues, and contribute to make this package better.

### How to Contribute

We welcome contributions from the community. If you encounter issues or have suggestions, please open an issue or submit a pull request.

