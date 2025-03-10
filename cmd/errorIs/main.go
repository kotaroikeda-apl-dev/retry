package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrNetwork = errors.New("ネットワークエラー")
var ErrOther = errors.New("その他のエラー")

func doJob() error {
	if rand.Intn(3) == 0 {
		return ErrNetwork
	}
	if rand.Intn(3) == 1 {
		return ErrOther
	}
	fmt.Println("ジョブ成功")
	return nil
}

func retryJob(retries int, delay time.Duration) error {
	var err error
	for i := 0; i < retries; i++ {
		err = doJob()
		if err == nil {
			return nil
		}
		if errors.Is(err, ErrOther) {
			return fmt.Errorf("致命的エラー: %w", err) // すぐに終了
		}
		fmt.Printf("リトライ %d/%d: エラー: %s\n", i+1, retries, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("リトライ失敗: %w", err)
}

func main() {
	err := retryJob(3, 1*time.Second)
	if err != nil {
		fmt.Println("最終エラー:", err)
	}
}
