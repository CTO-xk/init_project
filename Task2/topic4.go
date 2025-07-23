package main

import (
	"fmt"
	"sync"
	"time"
)

// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
type Task func()

func RunTasks(task []Task) []time.Duration {
	var wg sync.WaitGroup
	result := make([]time.Duration, len(task))
	for i, _ := range task {
		wg.Add(1)
		go func(taskId int, taskFun Task) {
			defer wg.Done()
			start := time.Now()             // 记录开始时间
			taskFun()                       // 执行任务
			duration := time.Since((start)) // 计算执行时间
			result[taskId] = duration       // 存储执行时间
		}(i, task[i])
	}
	wg.Wait()     // 等待所有任务完成
	return result // 返回每个任务的执行时间
}

func main() {
	//创建任务列表
	tasks := []Task{
		func() { time.Sleep(2. * time.Second); println("Task 1 completed") },
		func() { time.Sleep(2. * time.Second); println("Task 1 completed") },
		func() { time.Sleep(2. * time.Second); println("Task 1 completed") },
	}
	//执行任务的结果
	result := RunTasks(tasks)
	for _, res := range result {
		fmt.Printf("Task executed in: %v\n", res)
	}
}
