package main

import (
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	resp *http.Response
	err  error
}

func main() {
	done := make(chan struct{})
	producer := checkStatus([]string{"https://google.com", "https://onet.pl", "https://sth.c", "https://medium.com", "https://youtube.com", "https://olaf.comm", "http://error.error"}, done)

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

func checkStatus(urls []string, interupt <-chan struct{}) <-chan Result {
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
