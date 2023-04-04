package main

import (
	"fmt"
	"net/http"
	"time"
)

// The "Done" pattern in Go is a concurrency pattern that allows you to signal the completion of a task or goroutine
// This pattern is particularly useful when you need to coordinate multiple goroutines or tasks, and ensure that they all complete before the program exits.
// By using the "Done" pattern, you can create more robust and reliable programs that are better able to handle complex concurrency scenarios.

type Result struct {
	resp *http.Response
	err  error
}

func RunDone() {
	checkStatusJobConsumer()
}

func checkStatusJobConsumer() {
	done := make(chan struct{})
	producer := checkStatusJob([]string{"https://google.com", "https://onet.pl", "https://sth.c", "https://medium.com", "https://youtube.com", "https://olaf.comm", "http://error.error"}, done)

	var errCnt int
	for result := range producer {
		if result.err != nil {
			fmt.Printf("Error occured: %v\n", result.err)
			errCnt++
			if errCnt >= 3 {
				fmt.Println("Finishing program due to too many errs..")
				close(done)
			}
		} else {
			fmt.Printf("result: %v\n", result.resp.Status)
		}
	}
}

func checkStatusJob(urls []string, interupt <-chan struct{}) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)
		for _, url := range urls {
			time.Sleep(5 * time.Second)

			select {
			case <-interupt:
				return

			default:
				resp, err := http.Get(url)
				results <- Result{resp, err}
			}
		}
	}()

	return results
}
