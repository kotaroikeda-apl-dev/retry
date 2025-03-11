package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const maxRetries = 3

func shouldRetry(statusCode int) bool {
	// 500系エラーのみリトライ対象
	return statusCode >= 500 && statusCode < 600
}

func fetchData(url string) ([]byte, error) {
	var err error
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					return body, nil
				}
			} else if !shouldRetry(resp.StatusCode) {
				// 400系エラーはリトライしない
				return nil, fmt.Errorf("リクエスト失敗: ステータスコード %d", resp.StatusCode)
			}
		}

		if resp != nil {
			resp.Body.Close()
		}

		// バックオフを適用（1秒 → 2秒 → 4秒）
		delay := time.Duration(1<<i) * time.Second
		fmt.Printf("リトライ %d/%d: エラー: %s, 次の試行まで %v\n", i+1, maxRetries, err, delay)
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("リクエスト失敗: %w", err)
}

func main() {
	url := "https://httpstat.us/500" // 500エラーを返すURL
	data, err := fetchData(url)
	if err != nil {
		fmt.Println("最終エラー:", err)
	} else {
		fmt.Println("取得データ:", string(data))
	}
}
