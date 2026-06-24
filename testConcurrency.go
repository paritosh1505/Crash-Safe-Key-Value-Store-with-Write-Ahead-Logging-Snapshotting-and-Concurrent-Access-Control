package main

import (
	"fmt"
	"sync"
)

func incrementCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	count++
	fmt.Println(count)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go incrementCounter(&wg)
	go incrementCounter(&wg)
	wg.Wait()

}
