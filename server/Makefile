.PHONY: help build run test test-coverage lint fmt clean dev migrate-up migrate-down

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## バイナリをビルド
	go build -o bin/api ./cmd/api

run: ## サーバーを起動
	go run ./cmd/api

dev: ## 開発モード（ホットリロード）で起動
	air

test: ## テストを実行
	go test -v ./...

test-coverage: ## カバレッジ付きでテストを実行
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## リンターを実行
	golangci-lint run

fmt: ## コードをフォーマット
	gofumpt -l -w .
	goimports -w .

clean: ## ビルド成果物を削除
	rm -rf bin/ tmp/ coverage.out coverage.html

migrate-up: ## データベースマイグレーションを実行
	goose -dir migrations sqlite3 ./data.db up

migrate-down: ## データベースマイグレーションをロールバック
	goose -dir migrations sqlite3 ./data.db down

migrate-status: ## マイグレーションステータスを表示
	goose -dir migrations sqlite3 ./data.db status

.DEFAULT_GOAL := help

