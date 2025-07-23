package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func main() {
	var (
		count int64
		wg    sync.WaitGroup
	)
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// 使用原子操作递增计数器
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", count)
}
