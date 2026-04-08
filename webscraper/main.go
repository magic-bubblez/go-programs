package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Size       int64
	Duration   time.Duration
	Err        error
}

func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{URL: url, Err: err, Duration: duration}
			continue
		}
		results <- Result{
			URL:        url,
			StatusCode: resp.StatusCode,
			Size:       resp.ContentLength,
			Duration:   duration,
		}
		resp.Body.Close()
	}
}

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: webscraper <url1> <url2> ...")
		os.Exit(1)
	}

	//the Worker Pool pattern
	numw := 5
	if len(urls) < numw {
		numw = len(urls)
	}
	jobs := make(chan string, len(urls))    // wrokers pull url from here
	results := make(chan Result, len(urls)) // main pulls results from here

	var wg sync.WaitGroup
	for i := 0; i < numw; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}
	//could've feeded the urls before spawning the workers however standard pattern is
	//to always start workers first and then feed them (safe order regardless of buffered or unbuffered)
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	//collect and print all results
	for r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "  ERROR: %v (%v)\n\n", r.Err, r.Duration)
			continue
		}
		fmt.Printf("URL: %s\n", r.URL)
		fmt.Printf("  Status: %d | Size: %d bytes | Time: %v\n\n", r.StatusCode, r.Size, r.Duration)
	}
}
