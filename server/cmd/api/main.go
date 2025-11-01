package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ano333333/llm-time-manager/server/internal/database"
)

func main() {
	log.Println("LLM時間管理ツール - Server starting...")

	// データベースパスの設定
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/llm-time-manager.db"
	}

	// データベース接続の初期化
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	log.Printf("Database opened: %s", dbPath)

	// マイグレーションの実行
	migrationsDir := "./migrations"
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// マイグレーションバージョンを確認
	version, err := database.GetMigrationVersion(db)
	if err != nil {
		log.Printf("warning: failed to get migration version: %v", err)
	} else {
		log.Printf("Migration version: %d", version)
	}

	// TODO: サーバー設定の読み込み
	// TODO: HTTPサーバーの起動

	if _, err := fmt.Fprintln(os.Stdout, "Server is ready"); err != nil {
		log.Printf("failed to write to stdout: %v", err)
	}
}
