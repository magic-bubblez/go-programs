package main

import (
	"fmt"
	"sync"
)

var counter int
var wg sync.WaitGroup

func inc() {
	for i := 0; i < 1000000; i++ {
		counter++
	}
	wg.Done()
}
func main() {
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go inc()
	}
	wg.Wait()
	fmt.Println(counter)
}

// race condition - two goroutines accessing the same variable without synchronization
//both are constatnyl overwriting eachother's work. to fix this one can use mutex.
