package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func doJob() error {
	if rand.Intn(2) == 0 {
		return errors.New("ジョブ失敗")
	}
	fmt.Println("ジョブ成功")
	return nil
}

func retryWithBackoff(retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		err = doJob()
		if err == nil {
			return nil
		}
		delay := time.Duration(math.Pow(2, float64(i))) * time.Second
		fmt.Printf("リトライ %d/%d: エラー: %s, 次の試行まで %v\n", i+1, retries, err, delay)
		time.Sleep(delay)
	}
	return fmt.Errorf("リトライ失敗: %w", err)
}

func main() {
	err := retryWithBackoff(5)
	if err != nil {
		fmt.Println("最終エラー:", err)
	}
}
