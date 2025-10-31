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
	defer db.Close()

	log.Printf("Database opened: %s", dbPath)

	// マイグレーションの実行
	migrationsPath := "file://./migrations"
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// マイグレーションバージョンを確認
	version, dirty, err := database.GetMigrationVersion(db, migrationsPath)
	if err != nil {
		log.Printf("warning: failed to get migration version: %v", err)
	} else {
		log.Printf("Migration version: %d (dirty: %v)", version, dirty)
	}

	// TODO: サーバー設定の読み込み
	// TODO: HTTPサーバーの起動

	if _, err := fmt.Fprintln(os.Stdout, "Server is ready"); err != nil {
		log.Printf("failed to write to stdout: %v", err)
	}
}
