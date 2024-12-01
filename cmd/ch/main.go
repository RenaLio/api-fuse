package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan string)
	go func() {
		for i := range ch {
			fmt.Println(i)
		}
	}()
	wg := sync.WaitGroup{}
	for i := range 15 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- fmt.Sprintf("hello %d", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("DONE")
	time.Sleep(time.Second)
	close(ch)

}
