package main

import "fmt"

//题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func doubleSliceElements(slice *[]int) {
	nums := *slice
	for i := range nums {
		nums[i] *= 2
	}
}
func main() {
	nums := []int{1, 3, 5, 6, 7}
	doubleSliceElements(&nums)
	fmt.Println("The slice after doubling each element is:", nums)
}
