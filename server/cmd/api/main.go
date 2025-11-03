package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	setuphandlers "github.com/ano333333/llm-time-manager/server/cmd/api/setup"
	"github.com/ano333333/llm-time-manager/server/internal/database"
	"github.com/joho/godotenv"
)

const defaultDBPath = "./data/dev.db"
const defaultPort = 8080

func main() {
	log.Println("LLM時間管理ツール - Server starting...")

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("failed to load .env file: %v", err)
		}
	}

	// データベースパスの設定
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
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

	if _, err := fmt.Fprintln(os.Stdout, "Server is ready"); err != nil {
		log.Printf("failed to write to stdout: %v", err)
	}

	// PORTの取得
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = defaultPort
	}
	log.Printf("Server will be running on port %d", port)

	mux := setuphandlers.SetupHandlers(db)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	server.ListenAndServe()
}
