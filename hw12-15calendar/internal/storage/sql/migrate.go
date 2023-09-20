package sqlstorage

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	// Add Postgres driver for the "github.com/golang-migrate/migrate/v4" package.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func newMigrate(conf *Conf) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, fmt.Errorf("driver from fs.FS: %w", err)
	}

	//nolint:nosprintfhostport
	dbConfStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dbConfStr)
	if err != nil {
		return nil, fmt.Errorf("read migrations and connect to db: %w", err)
	}

	return m, nil
}

func MigrationsUp(conf *Conf) error {
	m, err := newMigrate(conf)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Warn("Migrations up: " + err.Error())
		} else {
			return fmt.Errorf("migrate up: %w", err)
		}
	}
	errSrc, errDB := m.Close()
	if errDB != nil {
		return fmt.Errorf("close db: %w", errDB)
	}
	if errSrc != nil {
		return fmt.Errorf("close src: %w", errSrc)
	}

	return nil
}

func MigrationsDown(conf *Conf) error {
	m, err := newMigrate(conf)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Warn("Migrations down: " + err.Error())
		} else {
			return fmt.Errorf("migrate up: %w", err)
		}
	}
	errSrc, errDB := m.Close()
	if errDB != nil {
		return fmt.Errorf("close db: %w", errDB)
	}
	if errSrc != nil {
		return fmt.Errorf("close src: %w", errSrc)
	}

	return nil
}
