package main

import (
	"fmt"
	"sync"
)

// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数
func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	do := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			<-ch //等待通道信号
			fmt.Printf("偶数: %d\n", i)
			ch <- struct{}{}
		}
		do <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 9; i += 2 {
			ch <- struct{}{} // 发送信号，表示可以处理奇数
			<-ch
			fmt.Printf("奇数: %d\n", i)
		}

	}()
	<-do
	wg.Wait()
	close(ch)
	close(do)
	fmt.Println("所有数字处理完毕")
}
