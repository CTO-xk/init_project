package main

import (
	"fmt"
	"sync"
)

// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
func main() {
	ch := make(chan int, 50)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			num := <-ch
			fmt.Printf("Received: %d\n", num)
		}
	}()
	wg.Wait()
}
