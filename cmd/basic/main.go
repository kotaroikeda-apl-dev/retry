package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func doJob() error {
	// 50%の確率でエラーを発生させる
	if rand.Intn(2) == 0 {
		return errors.New("ジョブ失敗")
	}
	fmt.Println("ジョブ成功")
	return nil
}

func retryJob(retries int, delay time.Duration) error {
	var err error
	for i := 0; i < retries; i++ {
		err = doJob()
		if err == nil {
			return nil // 成功したら終了
		}
		fmt.Printf("リトライ %d/%d: エラー: %s\n", i+1, retries, err)
		time.Sleep(delay) // リトライ間隔
	}
	return fmt.Errorf("リトライ失敗: %w", err)
}

func main() {
	err := retryJob(3, 1*time.Second)
	if err != nil {
		fmt.Println("最終エラー:", err)
	}
}
