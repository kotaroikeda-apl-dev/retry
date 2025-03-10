package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxRetries = 3

func doJob(id int) error {
	if rand.Intn(2) == 0 {
		return errors.New("ジョブ失敗")
	}
	fmt.Printf("Worker %d: ジョブ成功\n", id)
	return nil
}

func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	for job := range jobs {
		for i := 0; i < maxRetries; i++ {
			err = doJob(id)
			if err == nil {
				results <- fmt.Sprintf("Worker %d: ジョブ %d 成功", id, job)
				break
			}
			fmt.Printf("Worker %d: ジョブ %d リトライ %d/%d\n", id, job, i+1, maxRetries)
			time.Sleep(1 * time.Second)
		}

		if err != nil {
			results <- fmt.Sprintf("Worker %d: ジョブ %d 失敗 (リトライ上限超過)", id, job)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	jobs := make(chan int, 10)
	results := make(chan string, 10)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println(result)
	}
}
