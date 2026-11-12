package db

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"job-portal/internal/db/migrations"
)

func RunMigrations() error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	db, err := goose.OpenDBWithDriver("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("open db for goose: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
