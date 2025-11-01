package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

// RunMigrations はマイグレーションを実行する
// db: データベース接続
// migrationsDir: マイグレーションファイルのディレクトリパス（例: "./migrations"）
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// ディレクトリの存在確認
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migration directory not found: %w", err)
	}

	// gooseのダイアレクトを設定
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// マイグレーションを実行（最新バージョンまで）
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetMigrationVersion は現在のマイグレーションバージョンを取得する
func GetMigrationVersion(db *sql.DB) (int64, error) {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return 0, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return 0, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, nil
}

// DownMigration は1つ前のバージョンにロールバックする
func DownMigration(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// ResetMigrations はすべてのマイグレーションをロールバックする
func ResetMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Reset(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	return nil
}
