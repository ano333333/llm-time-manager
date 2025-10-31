package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations はマイグレーションを実行する
// dbPath: SQLiteデータベースファイルのパス
// migrationsPath: マイグレーションファイルのディレクトリパス（例: "file://./migrations"）
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// SQLite3ドライバーインスタンスを作成
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// マイグレーションインスタンスを作成
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// マイグレーションを実行（最新バージョンまで）
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetMigrationVersion は現在のマイグレーションバージョンを取得する
func GetMigrationVersion(db *sql.DB, migrationsPath string) (uint, bool, error) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}
