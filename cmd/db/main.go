package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL ドライバ
)

const maxRetries = 3

func executeQuery(ctx context.Context, db *sql.DB) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = db.ExecContext(ctx, "UPDATE users SET balance = balance - 100 WHERE id = 1")
		if err == nil {
			fmt.Println("クエリ成功")
			return nil
		}

		// 致命的なエラーなら即終了（例：デッドロック、接続喪失）
		if isFatalError(err) {
			return fmt.Errorf("即終了: %w", err)
		}

		fmt.Printf("リトライ %d/%d: エラー: %s\n", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * time.Second) // 1秒→2秒→3秒
	}
	return fmt.Errorf("クエリ失敗: %w", err)
}

// MySQL の致命的なエラーを判定
func isFatalError(err error) bool {
	// 例: デッドロックエラー (`Error 1213: Deadlock found when trying to get lock`)
	// 例: 接続喪失 (`Error 2006: MySQL server has gone away`)
	if err != nil {
		if err.Error() == "致命的エラー" || // カスタムエラー
			sql.ErrConnDone == err || // 接続が閉じられた
			sql.ErrTxDone == err { // トランザクションが終了
			return true
		}
	}
	return false
}

func main() {
	// MySQL の接続情報（適宜変更）
	dsn := "testuser:testpassword@tcp(127.0.0.1:3306)/testdb?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = executeQuery(ctx, db)
	if err != nil {
		log.Println("最終エラー:", err)
	}
}
