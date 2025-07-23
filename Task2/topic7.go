package main

import (
	"fmt"
	"sync"
)

// 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，
// 并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来
func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch) // 关闭通道，表示不再发送数据
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("Received: %d\n", num)
		}
	}()
	wg.Wait()

}
