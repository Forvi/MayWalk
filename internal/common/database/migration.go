package database

import (
	"errors"
	"log/slog"

	"github.com/Forvi/MayWalk/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrate(dbURL string, logger *slog.Logger) error {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("migration has not changed")
			return nil
		}

		return err
	}

	logger.Info("migrations successfully configured")
	return nil
}
