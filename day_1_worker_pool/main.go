package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(time.Second)
		results <- j * 2

	}

}

func main() {
	jobs := make(chan int, 3)
	results := make(chan int, 5)

	var wg sync.WaitGroup
	for i := 1; i < 4; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}
	close(jobs)
	go func() {
		for r := range results {
			fmt.Printf("主协程收到结果: %d\n", r)
		}
	}()

	wg.Wait()
	close(results)
}
