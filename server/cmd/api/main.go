package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("LLM時間管理ツール - Server starting...")

	// TODO: サーバー設定の読み込み
	// TODO: データベース接続の初期化
	// TODO: HTTPサーバーの起動

	if _, err := fmt.Fprintln(os.Stdout, "Server is ready"); err != nil {
		log.Printf("failed to write to stdout: %v", err)
	}
}
